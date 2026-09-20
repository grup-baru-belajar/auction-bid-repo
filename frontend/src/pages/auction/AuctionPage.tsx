import { useState, useEffect, useRef } from "react";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import {
  fetchAuctions,
  createAuction,
} from "../../features/auctions/auctionsSlice";
import { uploadImage, deleteImage } from "../../services/api/uploadApi";
import AuctionCard from "../../components/auction/AuctionCard";
import Pagination from "../../components/common/Pagination";

const LIMIT = 10;

const AuctionPage = () => {
  const dispatch = useAppDispatch();
  const {
    list: auctions,
    pagination,
    loading,
    error,
  } = useAppSelector((state) => state.auctions);

  const { user } = useAppSelector((state) => state.auth);
  const isAdmin = user?.role === "ADMIN";

  const [currentPage, setCurrentPage] = useState(1);
  const [isCompletedFilter, setIsCompletedFilter] = useState<
    boolean | undefined
  >(undefined);

  const [showModal, setShowModal] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [form, setForm] = useState({
    auctionName: "",
    description: "",
    startingPrice: "",
    endTime: "",
  });
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string>("");
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [formError, setFormError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadedImageUrl, setUploadedImageUrl] = useState<string | null>(null);
  const [uploadedImagePublicId, setUploadedImagePublicId] = useState<string | null>(null);

  useEffect(() => {
    dispatch(
      fetchAuctions({
        page: currentPage,
        limit: LIMIT,
        isCompleted: isCompletedFilter,
      }),
    );
  }, [currentPage, isCompletedFilter, dispatch]);

  useEffect(() => {
    return () => {
      if (imagePreview) URL.revokeObjectURL(imagePreview);
    };
  }, [imagePreview]);

  const handleFilterChange = (value: string) => {
    if (value === "all") setIsCompletedFilter(undefined);
    else setIsCompletedFilter(value === "completed");
    setCurrentPage(1);
  };

  const ALLOWED_IMAGE_TYPES = ["image/png", "image/jpeg", "image/webp"];

  const handleFile = (file: File) => {
    if (!ALLOWED_IMAGE_TYPES.includes(file.type)) {
      setFormError("Format gambar harus PNG, JPG, atau WEBP");
      return;
    }
    if (file.size > 1 * 1024 * 1024) {
      setFormError("Image size must be less than 1MB");
      return;
    }
    setFormError("");
    if (uploadedImagePublicId) {
      deleteImage(uploadedImagePublicId).catch(() => {});
    }
    setImageFile(file);
    setImagePreview(URL.createObjectURL(file));
    setUploadedImageUrl(null);
    setUploadedImagePublicId(null);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => setIsDragging(false);

  const resetForm = () => {
    setForm({
      auctionName: "",
      description: "",
      startingPrice: "",
      endTime: "",
    });
    setImageFile(null);
    setImagePreview("");
    setFormError("");
    setUploading(false);
    setUploadedImageUrl(null);
    setUploadedImagePublicId(null);
  };

  const handleValidate = () => {
    if (!form.auctionName.trim() || !form.description.trim() || !form.endTime) {
      setFormError("All fields are required");
      return;
    }

    const price = Number(form.startingPrice);

    if (!form.startingPrice || isNaN(price) || price <= 0) {
      setFormError("Starting price must be a number greater than 0");
      return;
    }

    const endDate = new Date(form.endTime);

    if (isNaN(endDate.getTime()) || endDate <= new Date()) {
      setFormError("End time must be in the future");
      return;
    }

    if (!imageFile) {
      setFormError("Image is required");
      return;
    }

    setFormError("");
    setShowModal(false);
    setShowConfirm(true);
  };

  const handleSubmit = async () => {
    if (!imageFile) return;
    setSubmitting(true);
    setFormError("");
    let imageLink = "";
    try {
      if (uploadedImageUrl) {
        imageLink = uploadedImageUrl;
      } else {
        setUploading(true);
        const result = await uploadImage(imageFile);
        imageLink = result.url;
        setUploadedImageUrl(result.url);
        setUploadedImagePublicId(result.publicId);
        setUploading(false);
      }
      await dispatch(
        createAuction({
          auctionName: form.auctionName,
          description: form.description,
          imageLink,
          startingPrice: Number(form.startingPrice),
          endTime: new Date(form.endTime).toISOString(),
        }),
      ).unwrap();
      setShowConfirm(false);
      resetForm();
      dispatch(
        fetchAuctions({
          page: currentPage,
          limit: LIMIT,
          isCompleted: isCompletedFilter,
        }),
      );
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { message?: string } } };
      const serverMsg = axiosErr?.response?.data?.message;
      if (uploadedImagePublicId) {
        deleteImage(uploadedImagePublicId).catch(() => {});
      }
      setFormError(
        serverMsg || "Failed to create the auction. Please try again.",
      );
      setShowConfirm(false);
      setShowModal(true);
    } finally {
      setSubmitting(false);
      setUploading(false);
    }
  };

  const startItem =
    pagination && auctions.length > 0 ? (currentPage - 1) * LIMIT + 1 : 0;
  const endItem = pagination ? (currentPage - 1) * LIMIT + auctions.length : 0;

  return (
    <div>
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-4">
          <p className="text-sm text-gray-500">
            Menampilkan <span className="text-black">{startItem}</span>-
            <span className="text-black">{endItem}</span> dari total{" "}
            <span className="text-black">{pagination?.total ?? 0}</span> Auction
          </p>

          <div className="flex items-center gap-3">
            {isAdmin && (
              <button
                onClick={() => setShowModal(true)}
                className="px-4 py-2 bg-[#1A4B69] hover:bg-[#12364c] text-white text-sm font-semibold rounded-lg transition-colors"
              >
                + Add Auction
              </button>
            )}
            <select
              onChange={(e) => handleFilterChange(e.target.value)}
              className="border border-gray-300 rounded-lg px-4 py-2 text-sm bg-white"
            >
              <option value="all">All</option>
              <option value="active">Active</option>
              <option value="completed">Completed</option>
            </select>
          </div>
        </div>

        {loading && (
          <p className="text-center text-gray-500 py-10">Loading Data...</p>
        )}
        {error && <p className="text-center text-red-500 py-10">{error}</p>}

        {!loading && !error && auctions.length === 0 && (
          <p className="text-center text-gray-500 py-10">No Auction Found.</p>
        )}

        {!loading && !error && auctions.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
              {auctions.map((auction) => (
                <AuctionCard key={auction.id} auction={auction} />
              ))}
            </div>

            <Pagination
              currentPage={currentPage}
              totalPages={pagination?.totalPages ?? 1}
              onPageChange={setCurrentPage}
            />
          </>
        )}
      </div>

      {showModal && (
        <div
          className="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
          onClick={() => {
            setShowModal(false);
            resetForm();
          }}
        >
          <div
            className="bg-white rounded-lg p-6 w-full max-w-md mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-bold mb-4">Add Auction</h2>
            <div className="space-y-3">
              <input
                type="text"
                placeholder="Auction Name"
                value={form.auctionName}
                onChange={(e) =>
                  setForm({ ...form, auctionName: e.target.value })
                }
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
              />
              <textarea
                placeholder="Description"
                value={form.description}
                onChange={(e) =>
                  setForm({ ...form, description: e.target.value })
                }
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm h-20 resize-none"
              />

              <div
                onDrop={handleDrop}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onClick={() => fileInputRef.current?.click()}
                className={`border-2 border-dashed rounded-lg p-4 text-center cursor-pointer transition-colors ${
                  isDragging
                    ? "border-blue-500 bg-blue-50"
                    : "border-gray-300 hover:border-gray-400"
                }`}
              >
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (file) handleFile(file);
                  }}
                />
                {imagePreview ? (
                  <div className="relative">
                    <img
                      src={imagePreview}
                      alt="Preview"
                      className="max-h-40 mx-auto rounded-lg object-contain"
                    />
                    <button
                      type="button"
                      onClick={(e) => {
                        e.stopPropagation();
                        if (uploadedImagePublicId) {
                          deleteImage(uploadedImagePublicId).catch(() => {});
                        }
                        setImageFile(null);
                        setImagePreview("");
                        setUploadedImageUrl(null);
                        setUploadedImagePublicId(null);
                        if (fileInputRef.current)
                          fileInputRef.current.value = "";
                      }}
                      className="absolute top-1 right-1 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center text-xs hover:bg-red-600"
                    >
                      x
                    </button>
                  </div>
                ) : (
                  <div className="py-4">
                    <svg
                      className="w-8 h-8 mx-auto text-gray-400 mb-2"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="1.5"
                    >
                      <path
                        d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                    <p className="text-sm text-gray-500">
                      Drag & drop gambar atau klik untuk browse
                    </p>
                    <p className="text-xs text-gray-400 mt-1">
                      PNG, JPG, WEBP (max 1MB)
                    </p>
                  </div>
                )}
              </div>

              <input
                type="number"
                placeholder="Starting Price"
                value={form.startingPrice}
                onChange={(e) =>
                  setForm({ ...form, startingPrice: e.target.value })
                }
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
              />
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  End Time
                </label>
                <input
                  type="datetime-local"
                  value={form.endTime}
                  onChange={(e) =>
                    setForm({ ...form, endTime: e.target.value })
                  }
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
                />
                <p className="text-xs text-gray-400 mt-1">
                  Start time will be set automatically when auction is created
                </p>
              </div>
              {formError && <p className="text-xs text-red-500">{formError}</p>}
            </div>
            <div className="flex justify-end gap-2 mt-5">
              <button
                onClick={() => {
                  setShowModal(false);
                  resetForm();
                }}
                className="px-4 py-2 text-sm text-gray-600 hover:text-gray-800"
              >
                Cancel
              </button>
              <button
                onClick={handleValidate}
                disabled={submitting}
                className="px-4 py-2 bg-[#1A4B69] hover:bg-[#12364c] text-white text-sm font-semibold rounded-lg disabled:opacity-60"
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}

      {showConfirm && (
        <div
          className="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
          onClick={() => {
            if (!submitting) {
              setShowConfirm(false);
              setShowModal(true);
            }
          }}
        >
          <div
            className="bg-white rounded-lg p-6 w-full max-w-sm mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-bold mb-2">Create Auction</h2>
            {submitting ? (
              <div className="mb-5 space-y-2">
                <div className="flex items-center gap-2 text-sm">
                  <div className={`w-4 h-4 border-2 rounded-full ${uploading ? "border-[#1A4B69] border-t-transparent animate-spin" : "border-green-500"}`} />
                  <span className={uploading ? "text-[#1A4B69] font-medium" : "text-green-600"}>
                    1. Uploading image...
                  </span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <div className={`w-4 h-4 border-2 rounded-full ${!uploading ? "border-[#1A4B69] border-t-transparent animate-spin" : "border-gray-300"}`} />
                  <span className={!uploading ? "text-[#1A4B69] font-medium" : "text-gray-400"}>
                    2. Creating auction...
                  </span>
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-600 mb-5">
                Are you sure you want to create this auction?
              </p>
            )}
            <div className="flex justify-end gap-2">
              <button
                onClick={() => {
                  setShowConfirm(false);
                  setShowModal(true);
                }}
                disabled={submitting}
                className="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 disabled:opacity-60"
              >
                Cancel
              </button>
              <button
                onClick={handleSubmit}
                disabled={submitting || uploading}
                className="px-4 py-2 bg-[#1A4B69] hover:bg-[#12364c] text-white text-sm font-semibold rounded-lg disabled:opacity-60"
              >
                {uploading ? "Uploading..." : submitting ? "Creating..." : "Create"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AuctionPage;
