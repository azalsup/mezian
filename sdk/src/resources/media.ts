/**
 * Media resource — image upload and YouTube video addition.
 * All operations require an access token.
 */

import type { HttpClient } from "../utils/http.js";
import type { Media } from "../types/index.js";

export class MediaResource {
  constructor(private http: HttpClient) {}

  /**
   * Uploads an image for an ad.
   * Accepte JPEG, PNG, WebP (max configuré dans config.yaml).
   * L'image est redimensionnée en thumbnail automatiquement.
   *
   * @param adId  - Ad ID (numeric, not slug)
   * @param file  - Fichier image (File API du navigateur ou Blob)
   *
   * @example
   * const input = document.querySelector<HTMLInputElement>('input[type=file]')!;
   * const media = await sdk.media.uploadImage(ad.id, input.files![0]);
   */
  uploadImage(adId: number, file: File | Blob): Promise<Media> {
    const formData = new FormData();
    formData.append("file", file);
    return this.http.upload(`/ads/${adId}/media`, formData);
  }

  /**
   * Adds a YouTube link to an ad (useful for real estate).
   *
   * @param adId     - Ad ID
   * @param youtubeUrl - URL complète YouTube, ex: "https://www.youtube.com/watch?v=..."
   */
  addYouTube(adId: number, youtubeUrl: string): Promise<Media> {
    return this.http.post(`/ads/${adId}/media/youtube`, { url: youtubeUrl });
  }

  /**
   * Deletes media (image or YouTube).
   */
  delete(mediaId: number): Promise<void> {
    return this.http.delete(`/media/${mediaId}`);
  }

  /**
   * Sets a media item as the ad cover image.
   */
  setCover(mediaId: number): Promise<Media> {
    return this.http.put(`/media/${mediaId}/cover`);
  }

  /**
   * Updates the display order of a media item.
   */
  updateOrder(mediaId: number, sortOrder: number): Promise<Media> {
    return this.http.put(`/media/${mediaId}/order`, { sort_order: sortOrder });
  }
}
