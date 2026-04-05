/**
 * Ressource Auth — authentication par OTP (phone) ou password.
 *
 * Recommended flow (phone):
 *   1. sendOtp({ phone, channel: "whatsapp" })
 *   2. verifyOtp({ phone, code })  → AuthResponse
 *
 * Password flow :
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
   * Requests sending an OTP by SMS or WhatsApp.
   * Le code arrive dans les secondes qui suivent.
   */
  sendOtp(req: SendOtpRequest): Promise<{ message: string }> {
    return this.http.post("/auth/send-otp", req);
  }

  /**
   * Verifies an OTP code and returns tokens + user profile.
   * Automatically creates the account if the user is new.
   */
  verifyOtp(req: VerifyOtpRequest): Promise<AuthResponse> {
    return this.http.post("/auth/verify-otp", req);
  }

  /**
   * Login by identifier (phone or email) + password.
   */
  login(req: LoginRequest): Promise<AuthResponse> {
    return this.http.post("/auth/login", req);
  }

  /**
   * Register with optional password.
   * A verification OTP is sent automatically.
   */
  register(req: RegisterRequest): Promise<AuthResponse> {
    return this.http.post("/auth/register", req);
  }

  /**
   * Refreshes the access token from the refresh token.
   */
  refresh(req: RefreshTokenRequest): Promise<AuthResponse> {
    return this.http.post("/auth/refresh", req);
  }

  /**
   * Logs out the user (revokes the refresh token server-side).
   * Requiert un access token valide.
   */
  logout(refreshToken: string): Promise<{ message: string }> {
    return this.http.post("/auth/logout", { refresh_token: refreshToken });
  }

  /** Returns the authenticated user's profile. */
  getMe(): Promise<User> {
    return this.http.get("/auth/me");
  }

  /** Updates the authenticated user's profile. */
  updateMe(req: UpdateMeRequest): Promise<User> {
    return this.http.put("/auth/me", req);
  }
}
