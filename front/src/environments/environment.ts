export const environment = {
  apiBaseUrl: '/api/v1',

  auth: {
    /** Set to false to hide OTP option and force password-only auth */
    otpEnabled: true,
    /** 'phone' | 'email' — which identifier is shown first on login/register */
    primaryIdentifier: 'phone' as 'phone' | 'email',
    /** Phone number is mandatory on registration */
    phoneRequired: true,
    /** Email is mandatory on registration */
    emailRequired: false,
    /** Pre-filled country on the address form (ISO 3166-1 alpha-2) */
    defaultCountry: 'MA',
  },
};
