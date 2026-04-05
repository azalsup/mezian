/**
 * Ressource Médias — upload d'images et ajout de vidéos YouTube.
 * Toutes les opérations requièrent un access token.
 */

import type { HttpClient } from "../utils/http.js";
import type { Media } from "../types/index.js";

export class MediaResource {
  constructor(private http: HttpClient) {}

  /**
   * Upload une image pour une annonce.
   * Accepte JPEG, PNG, WebP (max configuré dans config.yaml).
   * L'image est redimensionnée en thumbnail automatiquement.
   *
   * @param adId  - ID de l'annonce (numérique, pas le slug)
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
   * Ajoute un lien YouTube à une annonce (utile pour l'immobilier).
   *
   * @param adId     - ID de l'annonce
   * @param youtubeUrl - URL complète YouTube, ex: "https://www.youtube.com/watch?v=..."
   */
  addYouTube(adId: number, youtubeUrl: string): Promise<Media> {
    return this.http.post(`/ads/${adId}/media/youtube`, { url: youtubeUrl });
  }

  /**
   * Supprime un média (image ou YouTube).
   */
  delete(mediaId: number): Promise<void> {
    return this.http.delete(`/media/${mediaId}`);
  }

  /**
   * Définit un média comme image de couverture de l'annonce.
   */
  setCover(mediaId: number): Promise<Media> {
    return this.http.put(`/media/${mediaId}/cover`);
  }

  /**
   * Met à jour l'ordre d'affichage d'un média.
   */
  updateOrder(mediaId: number, sortOrder: number): Promise<Media> {
    return this.http.put(`/media/${mediaId}/order`, { sort_order: sortOrder });
  }
}
