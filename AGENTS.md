# Agents.md: Hacking CYOA Development Framework

This document serves as the technical and narrative specification for the **Hacking Choose Your Own Adventure (CYOA)** project. It is designed to ground any coding or creative agent in the specific UI constraints, state logic, and atmospheric requirements of the game.

---

## 1. Project Core Concept

A terminal-themed, high-stakes hacking thriller. The player takes the role of **Kael Voss**, an elite freelance operator navigating a web of corporate espionage, medical conspiracies, and personal vendettas.

### Narrative Pillars

- **Information Asymmetry:** The player rarely has the full picture. Intel is a currency.
- **Grounded Cyberpunk:** Avoid "magic hacking." Use terms like *Zero-day*, *Lateral Movement*, *Exfiltration*, and *Social Engineering*.
- **Moral Complexity:** Antagonists are driven by utilitarian logic, not just greed.

---

### UI Structure

- **Main Viewport:** Monospace font. Dark theme. 
- **Typewriter Effect:** Text must render character-by-character to simulate a terminal.
- **Secondary Page (Dossier):** A passive, non-interactive information bank that updates live as the `gameState` changes.

---
## 2. UI Component Breakdown

### **A. Header (Status Bar)**

The top section tracks game state variables through visual meters.

- Detection Meter: Labeled detection. Uses a 5-segment block gauge (square glyphs) to track stealth.

- Temporal Meter: Labeled time remaining. A thin, horizontal progress bar that fills or depletes based on the mission deadline.

- Identifier Tag: A small, rounded "pill" tag in the top right (e.g., ghost) representing the user's alias or current state.

### **B. Content Area (The Narrative Log)**
Standard Text: Clean, monospaced paragraphs.

Information Density: Mixes atmospheric "flavor text" with critical mission data (e.g., $35k bounty, 07:00 deadline).

System Observation: A highlighted line used for critical insights. It features:

A vertical accent bar (blue).

A comment prefix //.

Differentiation in color (Cyan) to separate "meta-thought" from "terminal output."

### **C. Interaction Layer (Choice Menu)**
Structure: Large, full-width rounded rectangular buttons.

Format: [Letter] Action — Context/Consequence.

Example: [A] phish an employee — send a spoofed IT reset email. slow, clean.

Hierarchy: The letter is bracketed to suggest a keyboard input requirement.

### **D. Command Line Interface (CLI) Prompt**
User String: root@ghost:~$ followed by a blinking underscore cursor _. This grounds the user in a Linux-style terminal environment.
All actions are tracked here by actual commands. Expected resulting output from commands is also shown, and includes changes to the game state
ie. resources and flags - [detection + 2]

---

## 3. The Dossier Specification

The Dossier is divided into three sections to provide context without breaking immersion:


| Section       | Content                                    | Logic                                               |
| ------------- | ------------------------------------------ | --------------------------------------------------- |
| **Intel**     | Organization details and project leaks.    | Progressively disclosed based on files found.       |
| **Contacts**  | NPC names, roles, and trust levels.        | Updates when NPCs are mentioned or met.             |
| **Situation** | Kael’s "Internal Monologue" + Heat Status. | Reflects current danger levels and immediate goals. |


---

## 4. Character Profiles

### **Kael Voss (Protagonist)**

- **Backstory:** Former federal contractor blacklisted after whistleblowing on illegal surveillance tools.
- **Motivation:** Financial survival + Uncovering the network that erased their career.
- **Trait:** Pragmatic, cynical, but prone to making "personal" choices.

### **Director Sandra Okoye (Antagonist)**

- **Role:** Chief of Clinical Strategy at Meridian Dynamics.
- **Motivation:** Believes a revolutionary drug (ARGO) is "too good to fail" despite dangerous trial anomalies.
- **Conflict:** Outcome-based thinker who views the protagonist as a threat to public health progress.

---

## 5. Scripting & Scene Logic

### Choice Weighting

Every player choice should impact at least one of the three "Resources":

1. **Heat:** Higher heat leads to tighter windows and more aggressive NPC reactions.
2. **Time:** Some endings are locked if the player takes too long.
3. **Identities:** Social engineering covers are "one-time use" or carry high risk of burning.

### Data-Driven Scene Format

Agents generating new content should follow this JSON-schema-inspired structure:

```json
{
  "id": "meridian_vault_01",
  "text": "The ARGO directory sits behind a biometric wall. You have the stolen hash, but cracking it will spike the server temperature.",
  "choices": [
    {
      "label": "Brute Force Hash",
      "cost": { "heat": 25, "time": 15 },
      "target": "meridian_vault_success"
    },
    {
      "label": "Spoof Admin Heartbeat",
      "requirement": "has_admin_token",
      "cost": { "heat": 5, "time": 10 },
      "target": "meridian_vault_stealth"
    }
  ],
  "sidebarUpdate": {
    "newIntel": "Biometric encryption detected.",
    "status": "Pulse rising. You're close."
  }
}
```

---

## 6. Style Guide

- **Tone:** Clinical, urgent, and professional. 
- **Visuals:** Use `[ ]` for choices and `>`  for terminal prompts.
- **Colors:** * `#33ff33` (Classic Green) - Standard Operations.
  - `#ffb300` (Amber) - Caution/Moderate Heat.
  - `#ff3333` (Red) - Critical Detection/System Failure.

---

## 7. Seed-Based Persistence System

To allow players to "resume" or "share" a specific game state, the engine uses a **State Seed**.

### Seed Anatomy: `[Path]-[Resources]-[Flags]`

Example Seed: `1479-H25T180-QmFzZTY0`

1. **Path (1479):** The sequence of Scene IDs visited. This allows the "Back" button to function by recalculating the state up to the previous ID.
2. **Resources (H25T180):** - `H25`: Heat is at 25%.
  - `T180`: 180 minutes remaining.
3. **Flags (QmFzZTY0):** A Base64 encoded bitfield.
  - Bit 0: Marcus Webb Contacted.
  - Bit 1: ARGO File Stolen.
  - Bit 2: Phishing identity burned.

### Why use a Seed?

- **Instant Save/Load:** The player just needs to copy-paste a small string.
- **Verification:** The agent can "replay" the seed to ensure the player didn't cheat (e.g., Heat cannot be 0 if the Path includes a brute-force scene).
- **Procedural Content:** A specific seed can influence the "random" flavor text in the Sidebar, making every "run" feel unique but repeatable.

---
## 8. Tamper-Proof Seed System (Anti-Cheat)
To prevent players from manually editing their stats (e.g., lowering Heat), the seed must be cryptographically signed using a Checksum or a simple XOR cipher.

### Seed Format: `[EncodedData].[Signature]`
Example: `MTQ3OS1IMjVUMTgw.b7a3`

1. **Encoded Data:** A Base64 string containing the Path, Resources, and Flags.
2. **Signature (The Seal):** A 4-character hex code generated by hashing the Encoded Data with a "Salt" (a secret internal string like `kael_voss_2026`).

### Validation Logic
When a seed is entered, the engine must:
1. Separate the `Data` from the `Signature`.
2. Re-calculate the hash of the `Data` using the internal Salt.
3. Compare the new hash to the provided `Signature`.
4. **If they don't match:** Trigger a "Data Corruption" or "Integrity Breach" screen.

### Coding Instruction: The "Brute Force" Trap
If the engine detects a tampered seed, do not just show an error. Display a terminal message:
`[CRITICAL] SESSION INTEGRITY BREACH DETECTED.`
`[WARN] IP LOGGED. TRACE INITIALIZED.`
`[SYSTEM] RELOADING LAST SECURE SNAPSHOT...`
This keeps the player inside the fiction even when they try to break it.
---
9. Command Injection & Secret Overrides
The game should listen for a text input buffer for certain choices. If the player types a specific command instead of clicking a button, the engine checks for a "Secret Override."

Logic Flow
The Prompt: Below the choice buttons, there is a blinking cursor > _. if clicked

The Hijack: If the input matches a secretCommand defined in the scene metadata, the standard choices are ignored, and a "Hidden Branch" is triggered.

The Reward: Secret commands typically result in 0 Heat or unique Intel because the player "knew the shortcut."

Example Implementation
Standard Choice: [A] Brute force the admin password (Heat +20)

Secret Command: Player types sudo login --bypass -u admin

Outcome: "You exploit a known kernel"

---
## 10. Technical Architecture

### Global State Management

The game engine must maintain a central `gameState` object to ensure choices have long-term consequences.

```javascript
{
  heat: 0,              // 0-100 scale (Detection level)
  timeRemaining: 240,    // In minutes, decrements per action
  inventory: [],         // Collected files, keys, and artifacts
  intel: [],             // Discovered facts for the Sidebar
  contacts: [],          // NPCs encountered and their status
  activeIdentities: [],  // Available social engineering covers
  flags: {               // Boolean triggers for plot milestones
    discoveredArgo: false,
    burnedMarcus: false,
    metSandra: false
  }
}
```