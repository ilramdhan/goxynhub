"use client";

import { useState, useRef } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/api/client";
import type { ApiResponse, Media } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

async function getMedia(siteId: string, page = 1): Promise<{ data: Media[]; meta: any }> {
  const { data } = await apiClient.get<ApiResponse<Media[]>>("/admin/media", {
    params: { site_id: siteId, page, per_page: 24 },
  });
  return { data: data.data ?? [], meta: data.meta };
}

async function deleteMedia(id: string): Promise<void> {
  await apiClient.delete(`/admin/media/${id}`);
}

export default function MediaPage() {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [selectedMedia, setSelectedMedia] = useState<Media | null>(null);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  const { data, isLoading } = useQuery({
    queryKey: ["media", SITE_ID],
    queryFn: () => getMedia(SITE_ID),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteMedia,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["media", SITE_ID] });
      setSelectedMedia(null);
    },
  });

  const handleUpload = async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    setUploading(true);
    setUploadProgress(0);

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const formData = new FormData();
      formData.append("file", file);
      formData.append("site_id", SITE_ID);

      try {
        await apiClient.post("/admin/media/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
          onUploadProgress: (progressEvent) => {
            const progress = Math.round(
              ((i + (progressEvent.loaded / (progressEvent.total || 1))) / files.length) * 100
            );
            setUploadProgress(progress);
          },
        });
      } catch (error) {
        console.error("Upload failed:", error);
      }
    }

    setUploading(false);
    setUploadProgress(0);
    queryClient.invalidateQueries({ queryKey: ["media", SITE_ID] });
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    handleUpload(e.dataTransfer.files);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  const copyToClipboard = (url: string) => {
    navigator.clipboard.writeText(url);
  };

  return (
    <div className="p-8">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Media Library</h1>
          <p className="text-gray-600 mt-1">Upload and manage your images and files</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex border border-gray-300 rounded-lg overflow-hidden">
            <button
              onClick={() => setViewMode("grid")}
              className={`px-3 py-2 text-sm ${viewMode === "grid" ? "bg-primary text-white" : "text-gray-600 hover:bg-gray-50"}`}
            >
              ⊞ Grid
            </button>
            <button
              onClick={() => setViewMode("list")}
              className={`px-3 py-2 text-sm ${viewMode === "list" ? "bg-primary text-white" : "text-gray-600 hover:bg-gray-50"}`}
            >
              ☰ List
            </button>
          </div>
          <button
            onClick={() => fileInputRef.current?.click()}
            className="px-4 py-2 bg-primary text-white rounded-lg font-semibold text-sm hover:bg-primary/90 transition-colors"
          >
            📤 Upload Files
          </button>
        </div>
      </div>

      {/* Upload zone */}
      <div
        onDrop={handleDrop}
        onDragOver={(e) => e.preventDefault()}
        className="border-2 border-dashed border-gray-300 rounded-xl p-8 mb-6 text-center hover:border-primary transition-colors cursor-pointer"
        onClick={() => fileInputRef.current?.click()}
      >
        {uploading ? (
          <div>
            <div className="text-2xl mb-2">⏳</div>
            <p className="text-gray-600 mb-2">Uploading... {uploadProgress}%</p>
            <div className="w-full bg-gray-200 rounded-full h-2 max-w-xs mx-auto">
              <div
                className="bg-primary h-2 rounded-full transition-all"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        ) : (
          <div>
            <div className="text-4xl mb-2">📁</div>
            <p className="text-gray-600">
              <span className="text-primary font-medium">Click to upload</span> or drag and drop
            </p>
            <p className="text-sm text-gray-400 mt-1">
              PNG, JPG, GIF, WebP, SVG, MP4, PDF (max 10MB)
            </p>
          </div>
        )}
      </div>

      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept="image/*,video/mp4,application/pdf"
        className="hidden"
        onChange={(e) => handleUpload(e.target.files)}
      />

      {/* Media grid/list */}
      <div className="flex gap-6">
        <div className="flex-1">
          {isLoading ? (
            <div className="text-center py-12 text-gray-500">Loading media...</div>
          ) : data?.data.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-4xl mb-3">🖼️</div>
              <p className="text-gray-500">No media files yet. Upload your first file!</p>
            </div>
          ) : viewMode === "grid" ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
              {data?.data.map((media) => (
                <button
                  key={media.id}
                  onClick={() => setSelectedMedia(media)}
                  className={`relative aspect-square rounded-xl overflow-hidden border-2 transition-all ${
                    selectedMedia?.id === media.id
                      ? "border-primary shadow-lg"
                      : "border-gray-200 hover:border-gray-300"
                  }`}
                >
                  {media.type === "image" ? (
                    <img
                      src={media.public_url}
                      alt={media.alt_text || media.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full bg-gray-100 flex items-center justify-center">
                      <span className="text-3xl">
                        {media.type === "video" ? "🎥" : media.type === "document" ? "📄" : "📁"}
                      </span>
                    </div>
                  )}
                  <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white text-xs p-1 truncate">
                    {media.name}
                  </div>
                </button>
              ))}
            </div>
          ) : (
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">File</th>
                    <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Type</th>
                    <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Size</th>
                    <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Date</th>
                    <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data?.data.map((media) => (
                    <tr
                      key={media.id}
                      className="hover:bg-gray-50 cursor-pointer"
                      onClick={() => setSelectedMedia(media)}
                    >
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          {media.type === "image" ? (
                            <img
                              src={media.public_url}
                              alt={media.name}
                              className="w-10 h-10 object-cover rounded"
                            />
                          ) : (
                            <div className="w-10 h-10 bg-gray-100 rounded flex items-center justify-center text-xl">
                              {media.type === "video" ? "🎥" : "📄"}
                            </div>
                          )}
                          <span className="text-sm font-medium text-gray-900">{media.name}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-500">{media.mime_type}</td>
                      <td className="px-4 py-3 text-sm text-gray-500">{formatFileSize(media.file_size)}</td>
                      <td className="px-4 py-3 text-sm text-gray-500">
                        {new Date(media.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            copyToClipboard(media.public_url);
                          }}
                          className="text-xs text-primary hover:underline mr-3"
                        >
                          Copy URL
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            if (confirm("Delete this file?")) {
                              deleteMutation.mutate(media.id);
                            }
                          }}
                          className="text-xs text-red-600 hover:underline"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Media detail panel */}
        {selectedMedia && (
          <div className="w-72 flex-shrink-0 bg-white rounded-xl border border-gray-200 p-4 h-fit sticky top-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold text-gray-900 text-sm">File Details</h3>
              <button
                onClick={() => setSelectedMedia(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                ✕
              </button>
            </div>

            {selectedMedia.type === "image" && (
              <img
                src={selectedMedia.public_url}
                alt={selectedMedia.alt_text || selectedMedia.name}
                className="w-full rounded-lg mb-3 border border-gray-200"
              />
            )}

            <div className="space-y-2 text-sm">
              <div>
                <span className="text-gray-500">Name:</span>
                <span className="ml-2 font-medium text-gray-900">{selectedMedia.name}</span>
              </div>
              <div>
                <span className="text-gray-500">Type:</span>
                <span className="ml-2 text-gray-700">{selectedMedia.mime_type}</span>
              </div>
              <div>
                <span className="text-gray-500">Size:</span>
                <span className="ml-2 text-gray-700">{formatFileSize(selectedMedia.file_size)}</span>
              </div>
              {selectedMedia.width && selectedMedia.height && (
                <div>
                  <span className="text-gray-500">Dimensions:</span>
                  <span className="ml-2 text-gray-700">
                    {selectedMedia.width} × {selectedMedia.height}
                  </span>
                </div>
              )}
              <div>
                <span className="text-gray-500">Uploaded:</span>
                <span className="ml-2 text-gray-700">
                  {new Date(selectedMedia.created_at).toLocaleDateString()}
                </span>
              </div>
            </div>

            <div className="mt-4 space-y-2">
              <button
                onClick={() => copyToClipboard(selectedMedia.public_url)}
                className="w-full py-2 bg-primary/10 text-primary rounded-lg text-sm font-medium hover:bg-primary/20 transition-colors"
              >
                📋 Copy URL
              </button>
              <a
                href={selectedMedia.public_url}
                target="_blank"
                rel="noopener noreferrer"
                className="block w-full py-2 border border-gray-300 text-gray-700 rounded-lg text-sm font-medium text-center hover:bg-gray-50 transition-colors"
              >
                🔗 Open in New Tab
              </a>
              <button
                onClick={() => {
                  if (confirm("Delete this file?")) {
                    deleteMutation.mutate(selectedMedia.id);
                  }
                }}
                className="w-full py-2 border border-red-200 text-red-600 rounded-lg text-sm font-medium hover:bg-red-50 transition-colors"
              >
                🗑️ Delete File
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
