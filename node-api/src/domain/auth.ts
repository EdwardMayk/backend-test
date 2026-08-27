import crypto from 'crypto';

function base64urlDecode(str: string): string {
  let base64 = str.replace(/-/g, '+').replace(/_/g, '/');
  while (base64.length % 4) {
    base64 += '=';
  }
  return Buffer.from(base64, 'base64').toString('utf8');
}

export function verifyJWT(token: string): any {
  const secret = process.env.JWT_SECRET;
  if (!secret) {
    throw new Error('la variable de entorno JWT_SECRET no esta configurada');
  }

  const parts = token.split('.');
  if (parts.length !== 3) {
    throw new Error('estructura de token invalida');
  }

  const [headerB64, payloadB64, signatureB64] = parts;

  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(`${headerB64}.${payloadB64}`)
    .digest('base64url');

  if (signatureB64 !== expectedSignature) {
    throw new Error('verificacion de firma fallida');
  }

  const payload = JSON.parse(base64urlDecode(payloadB64));
  if (payload.exp && Date.now() >= payload.exp * 1000) {
    throw new Error('el token ha expirado');
  }

  return payload;
}
