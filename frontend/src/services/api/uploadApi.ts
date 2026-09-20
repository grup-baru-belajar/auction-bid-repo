import apiClient from "./apiClient";

export interface UploadResult {
  url: string;
  publicId: string;
}

export async function uploadImage(file: File): Promise<UploadResult> {
  const formData = new FormData();
  formData.append("image", file);

  const res = await apiClient.post("/v1/upload/image", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });

  const data = res.data;
  if (!data.success) {
    throw new Error(data.message || "Failed to upload image");
  }

  return {
    url: data.data.url as string,
    publicId: data.data.public_id as string,
  };
}

export async function deleteImage(publicId: string): Promise<void> {
  await apiClient.delete(`/v1/upload/image?public_id=${encodeURIComponent(publicId)}`);
}
