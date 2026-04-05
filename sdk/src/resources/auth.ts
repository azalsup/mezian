/**
 * Ressource Auth — authentification par OTP (téléphone) ou mot de passe.
 *
 * Flux recommandé (téléphone) :
 *   1. sendOtp({ phone, channel: "whatsapp" })
 *   2. verifyOtp({ phone, code })  → AuthResponse
 *
 * Flux mot de passe :
 *   1. login({ identifier: "+212...", password })  → AuthResponse
 */

import type { HttpClient } from "../utils/http.js";
import type {
  AuthResponse,
  LoginRequest,
  RefreshTokenRequest,
  RegisterRequest,
  SendOtpRequest,
  UpdateMeRequest,
  User,
  VerifyOtpRequest,
} from "../types/index.js";

export class AuthResource {
  constructor(private http: HttpClient) {}

  /**
   * Demande l'envoi d'un OTP par SMS ou WhatsApp.
   * Le code arrive dans les secondes qui suivent.
   */
  sendOtp(req: SendOtpRequest): Promise<{ message: string }> {
    return this.http.post("/auth/send-otp", req);
  }

  /**
   * Vérifie un code OTP et retourne les tokens + profil utilisateur.
   * Crée le compte automatiquement si l'utilisateur est nouveau.
   */
  verifyOtp(req: VerifyOtpRequest): Promise<AuthResponse> {
    return this.http.post("/auth/verify-otp", req);
  }

  /**
   * Connexion par identifiant (téléphone ou email) + mot de passe.
   */
  login(req: LoginRequest): Promise<AuthResponse> {
    return this.http.post("/auth/login", req);
  }

  /**
   * Inscription avec mot de passe (optionnel).
   * Un OTP de vérification est envoyé automatiquement.
   */
  register(req: RegisterRequest): Promise<AuthResponse> {
    return this.http.post("/auth/register", req);
  }

  /**
   * Renouvelle l'access token à partir du refresh token.
   */
  refresh(req: RefreshTokenRequest): Promise<AuthResponse> {
    return this.http.post("/auth/refresh", req);
  }

  /**
   * Déconnecte l'utilisateur (révoque le refresh token côté serveur).
   * Requiert un access token valide.
   */
  logout(refreshToken: string): Promise<{ message: string }> {
    return this.http.post("/auth/logout", { refresh_token: refreshToken });
  }

  /** Retourne le profil de l'utilisateur connecté. */
  getMe(): Promise<User> {
    return this.http.get("/auth/me");
  }

  /** Met à jour le profil de l'utilisateur connecté. */
  updateMe(req: UpdateMeRequest): Promise<User> {
    return this.http.put("/auth/me", req);
  }
}
