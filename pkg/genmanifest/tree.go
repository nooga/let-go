/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package genmanifest

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TreeManifestName is the sentinel file lgbgen writes at the root of the
// generated lowered tree (pkg/rt/core_go_lowered) as the LAST step of a
// successful generation. It lists a sha256 per generated file, so consumers
// can positively distinguish "complete tree" from "torn or partially written
// tree" instead of discovering the latter via cryptic build failures. The
// generator removes it first thing on reinstall, so a crash mid-install can
// never leave a valid sentinel over an inconsistent tree.
const TreeManifestName = ".lgbgen-tree.sum"

// hashFile returns the hex sha256 of one file's contents.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// treeFiles returns the sorted slash-relative paths of every regular file
// under dir, excluding the manifest itself.
func treeFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == TreeManifestName {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// WriteTreeManifest hashes every file under dir and writes the sentinel
// manifest ("<sha256>  <relpath>" per line, sorted). Call it only after the
// tree is fully written — its presence asserts completeness.
func WriteTreeManifest(dir string) error {
	files, err := treeFiles(dir)
	if err != nil {
		return err
	}
	var b strings.Builder
	for _, rel := range files {
		sum, err := hashFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "%s  %s\n", sum, rel)
	}
	return os.WriteFile(filepath.Join(dir, TreeManifestName), []byte(b.String()), 0644)
}

// CheckTreeManifest verifies the generated tree at dir against its sentinel
// manifest: the sentinel must exist, every listed file must be present with
// matching contents, and no unlisted files may have appeared. A nil error
// means the tree is exactly the one a successful generation wrote.
func CheckTreeManifest(dir string) error {
	f, err := os.Open(filepath.Join(dir, TreeManifestName))
	if os.IsNotExist(err) {
		return fmt.Errorf("no %s in %s — tree incomplete or never generated", TreeManifestName, dir)
	}
	if err != nil {
		return err
	}
	defer f.Close()

	want := map[string]string{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		sum, rel, ok := strings.Cut(line, "  ")
		if !ok {
			return fmt.Errorf("malformed %s line: %q", TreeManifestName, line)
		}
		want[rel] = sum
	}
	if err := sc.Err(); err != nil {
		return err
	}

	have, err := treeFiles(dir)
	if err != nil {
		return err
	}
	haveSet := map[string]bool{}
	for _, rel := range have {
		haveSet[rel] = true
		sum, ok := want[rel]
		if !ok {
			return fmt.Errorf("unlisted file %s in %s (stray or partial regeneration)", rel, dir)
		}
		got, err := hashFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		if got != sum {
			return fmt.Errorf("checksum mismatch for %s in %s", rel, dir)
		}
	}
	for rel := range want {
		if !haveSet[rel] {
			return fmt.Errorf("missing file %s in %s (torn tree)", rel, dir)
		}
	}
	return nil
}
