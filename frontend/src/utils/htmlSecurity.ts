// HTML sanitization function
export const sanitizeHTML = (html: string): string => {
  // Remove script tags and their content (case insensitive)
  let sanitized = html.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, "");

  // Remove dangerous event handlers (more comprehensive)
  sanitized = sanitized.replace(/\son\w+\s*=\s*["'][^"']*["']/gi, "");
  sanitized = sanitized.replace(/\son\w+\s*=\s*[^"'\s>]+/gi, "");

  // Remove javascript: URLs
  sanitized = sanitized.replace(/javascript:[^"'\s>]*/gi, "");

  // Remove data URLs that could contain scripts
  sanitized = sanitized.replace(/data:[^"']*base64[^"']*/gi, "");

  // Remove vbscript: URLs
  sanitized = sanitized.replace(/vbscript:[^"'\s>]*/gi, "");

  // Remove object, embed tags (but not iframe since we control it)
  sanitized = sanitized.replace(/<(object|embed)\b[^>]*>.*?<\/\1>/gi, "");
  
  // Remove dangerous tags
  sanitized = sanitized.replace(/<(applet|meta|link)\b[^>]*>.*?<\/\1>/gi, "");
  sanitized = sanitized.replace(/<(applet|meta|link)\b[^>]*\/>/gi, "");

  // Remove form elements that could be used for attacks
  sanitized = sanitized.replace(
    /<(form|input|textarea|button|select|option)\b[^>]*>.*?<\/\1>/gi,
    ""
  );
  sanitized = sanitized.replace(/<(form|input|textarea|button|select|option)\b[^>]*\/>/gi, "");

  // Wrap in a basic HTML structure for proper rendering
  return `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <style>
        body { font-family: system-ui, sans-serif; padding: 1rem; margin: 0; }
        * { max-width: 100%; }
      </style>
    </head>
    <body>
      ${sanitized}
    </body>
    </html>
  `;
};

// Check if HTML contains potentially dangerous content
export const hasSecurityRisk = (html: string): string | null => {
  if (/<script\b/i.test(html)) {
    return "Script tags detected";
  }
  if (/\son\w+\s*=/i.test(html)) {
    return "Event handlers detected";
  }
  if (/javascript:/i.test(html)) {
    return "JavaScript URLs detected";
  }
  if (/vbscript:/i.test(html)) {
    return "VBScript URLs detected";
  }
  if (/<(object|embed|applet|form|input|textarea|button|select|option)\b/i.test(html)) {
    return "Potentially dangerous HTML elements detected";
  }
  return null;
};
