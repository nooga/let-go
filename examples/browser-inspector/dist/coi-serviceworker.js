addEventListener('install', () => skipWaiting());
addEventListener('activate', e => e.waitUntil(clients.claim()));
addEventListener('fetch', e => {
  if (e.request.cache === 'only-if-cached' && e.request.mode !== 'same-origin') return;
  e.respondWith(fetch(e.request).then(r => {
    if (r.status === 0) return r;
    const h = new Headers(r.headers);
    // Pass server-set isolation headers through untouched. Overriding them
    // (the previous behavior) broke require-corp setups on dev servers
    // that already provide proper headers, by replacing them with the
    // credentialless variant Safari rejects — yielding pages that look
    // like they should be isolated but aren't.
    if (!h.has('Cross-Origin-Embedder-Policy')) {
      // require-corp is the broadest-compatible option: Safari, Firefox,
      // and Chrome all accept it; credentialless is Chrome-only.
      h.set('Cross-Origin-Embedder-Policy', 'require-corp');
    }
    if (!h.has('Cross-Origin-Opener-Policy')) {
      h.set('Cross-Origin-Opener-Policy', 'same-origin');
    }
    return new Response(r.body, {status: r.status, statusText: r.statusText, headers: h});
  }).catch(() => new Response(null, {status: 500})));
});
