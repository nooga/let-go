package wasm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBrowserInspectorShellHasDedicatedArtifactPanes(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "examples", "browser-inspector", "shell.html"))
	if err != nil {
		t.Fatalf("read shell template: %v", err)
	}
	got := string(data)
	for _, want := range []string{
		`id="artifact-result"`,
		`id="artifact-stdout"`,
		`id="artifact-stderr"`,
		`id="artifact-bytecode"`,
		`id="artifact-ir"`,
		`id="artifact-optimized-bytecode"`,
		`id="artifact-lowered-go"`,
		`const artifactTemplates = {`,
		`bytecode(name, payload) {`,
		`bytecodeRows(name, payload) {`,
		`const artifactDefs = {`,
		`frameArtifact: 'optimized_bytecode',`,
		`const artifactByFrameName = Object.fromEntries(`,
		`function formatNamedArtifact(name, payload)`,
		`const orderedArtifacts = ['result', 'stdout', 'stderr', 'bytecode', 'ir', 'optimizedBytecode', 'loweredGo'];`,
		`function createCell(source = '') {`,
		`const seedForms = splitTopLevelForms(defaultCellSource);`,
		`const cells = (seedForms.length ? seedForms : ['']).map((form) => createCell(form));`,
		`function renderActiveCellArtifacts() {`,
		`function setCellArtifactState(cell, name, status, content) {`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("shell template missing %q", want)
		}
	}
}

func TestBrowserInspectorShellHasAnalysisToggles(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "examples", "browser-inspector", "shell.html"))
	if err != nil {
		t.Fatalf("read shell template: %v", err)
	}
	got := string(data)
	for _, want := range []string{
		`<fieldset class="toggles">`,
		`<legend>Analysis</legend>`,
		`id="add-cell-btn"`,
		`id="clear-session-btn"`,
		`delBtn.title = 'Delete this cell';`,
		`id="cell-list"`,
		`id="toggle-bytecode"`,
		`id="toggle-ir"`,
		`id="toggle-optimized-bytecode"`,
		`id="toggle-lowered-go"`,
		`function selectedInspectOptions() {`,
		`function requestedArtifactNames(inspect) {`,
		`Notebook-style preview over <code>LetGoHost.request(...)</code> with one shared evaluation session.`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("shell template missing %q", want)
		}
	}
	for _, unwanted := range []string{
		`IR output is not wired yet.`,
		`Optimized bytecode output is not wired yet.`,
		`Lowered Go output is not wired yet.`,
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("shell template still contains stale unsupported placeholder %q", unwanted)
		}
	}
}

func TestHostEvalOutputUsesStreamedRequestFrames(t *testing.T) {
	got := AssembleHTML("// stub wasm_exec.js\n", "STUBWASM==", false, true, true)
	for _, want := range []string{
		`window.LetGoHost.request = function(req, onFrame)`,
		`postMessage({t:'request-frame', id: e.data.id, frame: frameJson});`,
		`if (e.data.t === 'request-frame') {`,
		`requestImpl = (req, onFrame) => new Promise((resolve, reject) => {`,
		`requestPending.set(id, { resolve, reject, timer, onFrame, frames: [] });`,
		`requestImpl = (req, onFrame) => {`,
		`window._lgRequestFrame = function(frameJson) {`,
		`return Promise.resolve({ ack, frames });`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("assembled host-eval output missing %q", want)
		}
	}
	for _, unwanted := range []string{
		`if (!resp || !resp.ok) {`,
		`const artifacts = resp.artifacts || {};`,
		`renderResultArtifact(resp);`,
		"op,\n          code: codeEl.value",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("assembled host-eval output still contains stale aggregated contract %q", unwanted)
		}
	}
}

func TestBrowserInspectorShellRendersStreamedFrames(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "examples", "browser-inspector", "shell.html"))
	if err != nil {
		t.Fatalf("read shell template: %v", err)
	}
	got := string(data)
	for _, want := range []string{
		`function inspectFor(op) {`,
		`function renderFrame(cell, frame) {`,
		`let lastOp = 'eval';`,
		`let activeRun = 0;`,
		`let activeCellID = null;`,
		`let session = ` + "`browser-inspector-${Date.now()}`" + `;`,
		`let reanalyzeTimer = null;`,
		`const REANALYZE_DEBOUNCE_MS = 250;`,
		`function renderCellList() {`,
		`function setActiveCell(id) {`,
		`function clearSessionState() {`,
		`function resetCellArtifacts(cell, op, inspect) {`,
		`const inspect = inspectFor(op);`,
		`const requested = new Set(requestedArtifactNames(inspect));`,
		`const received = new Set();`,
		`const resp = await window.LetGoHost.request({`,
		`frames.forEach((frame) => {`,
		`if (name && !received.has(name)) {`,
		`received.add(name);`,
		`appendCellArtifactSection(cell, name, frame.form_index || 0, rendered.status, rendered.text);`,
		`op: 'eval',`,
		`inspect,`,
		`}, (frame) => {`,
		`if (frame && (frame.kind === 'artifact' || frame.kind === 'err' || frame.kind === 'note')) {`,
		`renderFrame(cell, frame);`,
		`if (frame.kind === 'out') {`,
		`appendCellArtifactText(cell, 'stdout', frame.text || '', 'ok');`,
		`if (frame.kind === 'err') {`,
		`const artifactName = artifactByFrameName[frame.artifact];`,
		`appendCellArtifactSection(cell, artifactName, frame.form_index || 0, 'error',`,
		`appendCellArtifactText(cell, 'stderr', frame.text || 'request failed', 'error');`,
		`document.getElementById('add-cell-btn').addEventListener('click', addCell);`,
		`document.getElementById('clear-session-btn').addEventListener('click', clearSessionState);`,
		`toggleEls.bytecode.addEventListener('change', () => send(lastOp));`,
		`toggleEls.ir.addEventListener('change', () => send(lastOp));`,
		`toggleEls.optimizedBytecode.addEventListener('change', () => send(lastOp));`,
		`toggleEls.loweredGo.addEventListener('change', () => send(lastOp));`,
		`renderCellList();`,
		`editor.addEventListener('input', () => { cell.source = editor.value; autosize(editor); });`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("shell template missing %q", want)
		}
	}
	for _, unwanted := range []string{
		`if (frame.artifact === 'optimized_bytecode') {`,
		`if (frame.artifact === 'lowered_go') {`,
		`if (!resp.ok) {`,
		`const artifacts = resp.artifacts || {};`,
		`renderResultArtifact(resp);`,
		"op,\n          code: codeEl.value",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("shell template still contains stale aggregated contract %q", unwanted)
		}
	}
}

func TestBrowserInspectorBuildScriptDefaultsToSourceDrivenLG(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "examples", "browser-inspector", "Makefile"))
	if err != nil {
		t.Fatalf("read build Makefile: %v", err)
	}
	got := string(data)
	// With no LG override, the build compiles lg from source (go run ./),
	// not a checked-in ./lg binary.
	if !strings.Contains(got, `go run ./`) {
		t.Fatalf("build Makefile must default to a source-driven lg invocation")
	}
	if strings.Contains(got, `LG ?= $(ROOT)/lg`) || strings.Contains(got, `LG := $(ROOT)/lg`) {
		t.Fatalf("build Makefile still defaults to the checked-in ./lg binary")
	}
}
