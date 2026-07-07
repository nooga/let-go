// --- Default xterm.js shell ---
// The shell `lg -w` ships unless built with -w-shell none. It owns the
// presentation layer (the Terminal, its font/theme, the DOM) and binds to
// the runtime purely through window.LetGoHost — the same surface a client's
// own shell uses. Nothing in lg-host-core.js references xterm; everything
// terminal-specific lives here.
(function() {
  const status = document.getElementById('status');
  const termEl = document.getElementById('app');

  const term = new Terminal({
    fontFamily: '"IBM Plex Mono", "Menlo", "Consolas", monospace',
    fontSize: 14,
    theme: { background: '#0c0c0c', foreground: '#e8e6df', cursor: '#5ec4b6' },
    allowProposedApi: true,
    convertEol: true,
  });
  const fitAddon = new FitAddon.FitAddon();
  term.loadAddon(fitAddon);

  function showTerminal() {
    if (status) status.style.display = 'none';
    term.open(termEl);
    fitAddon.fit();
    term.focus();
  }

  window.addEventListener('resize', () => fitAddon.fit());

  // Bind once the core glue is wired. onOutput routes VM stdout to xterm;
  // in worker mode the keystroke/size path is live, so advertise the
  // initial size and forward xterm input + resizes through LetGoHost.
  window.LetGoHost.onReady((mode) => {
    showTerminal();
    window.LetGoHost.onOutput((s) => term.write(s));
    if (mode === 'worker') {
      window.LetGoHost.setSize(term.cols, term.rows);
      term.onResize(({cols, rows}) => window.LetGoHost.setSize(cols, rows));
      term.onData((data) => window.LetGoHost.sendInput(data));
    }
  });
})();
