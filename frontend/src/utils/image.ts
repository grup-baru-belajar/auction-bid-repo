import blockedPlaceholder from "../assets/react.svg";

export const IMAGE_BLOCKED_BY_NETWORK = "block_by_network";

export const NETWORK_ONLY_ERROR_CODES = ["NETWORK_BLOCKED"];

export const getUploadErrorCode = (err: unknown): string | undefined =>
  (err as { response?: { data?: { error?: { code?: string } } } })?.response
    ?.data?.error?.code;

export const isNetworkOnlyError = (err: unknown): boolean => {
  const code = getUploadErrorCode(err);
  return !!code && NETWORK_ONLY_ERROR_CODES.includes(code);
};

export const resolveImageUrl = (imageLink?: string | null): string =>
  !imageLink || imageLink === IMAGE_BLOCKED_BY_NETWORK
    ? blockedPlaceholder
    : imageLink;
