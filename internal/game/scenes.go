package game

func Scenes() map[string]Scene {
	return map[string]Scene{

		"intro": {
			Prompt:  "root@ghost:~$",
			Mystery: "// internal-legacy.meridian.net — why is this still running?",
			Text: `03:47 AM. Your terminal blinks against the dark.

Meridian Dynamics login portal. $35k for one file — internal trial
data, folder unknown. IT rotation at 07:00. Three hours, twelve
minutes.

Recon turned up two clean entry points. But there's a third you
didn't expect: internal-legacy.meridian.net. No SSL. Last
modified: 2009.

Why is it still live?`,
			TargetEntries: []DossierEntry{
				{"Meridian Dynamics", "Mid-tier pharma. 428 staff. HQ: Austin TX.", ""},
				{"Target file", "Internal trial data — folder path unknown.", ""},
				{"Client", "Anonymous. $35k on delivery. No questions asked.", ""},
			},
			SituationLines: []StatusLine{
				{"No alerts triggered.", "ok"},
				{"3h 12m before IT rotation.", "ok"},
			},
			Choices: []Choice{
				{
					Label:     "[A] phish an employee — spoofed IT reset. slow, clean.",
					Next:      "phish_setup",
					Type:      ChoiceSafe,
					HeatCost:  0,
					TimeCost:  15,
					Archetype: "ghost",
				},
				{
					Label:     "[B] exploit VPN client CVE — fast entry, noisy.",
					Next:      "vpn_exploit",
					Type:      ChoiceDanger,
					HeatCost:  2,
					TimeCost:  5,
					Archetype: "coder",
				},
				{
					Label:    "[C] probe the legacy subdomain — something feels off.",
					Next:     "legacy_probe",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 8,
				},
			},
		},

		"phish_setup": {
			Prompt:  "root@ghost:~/tools/phish$",
			Mystery: "// forwarded inbox — recipient unknown. could be IR.",
			Text: `LinkedIn returns a name: Dana Reyes, IT helpdesk. Active recently.

But she's been on paternity leave since last week. Her email is
being forwarded somewhere. You don't know who receives it.

Spoofed reset email is ready. "Mandatory password change — 30
minutes." Standard pretext. Question is whether anyone on the
other end is awake.`,
			ContactEntries: []DossierEntry{
				{"Dana Reyes", "IT helpdesk. On leave. Email forwarded — recipient unknown.", "al"},
			},
			SituationLines: []StatusLine{
				{"No alerts triggered.", "ok"},
				{"Phish ready. Recipient identity uncertain.", "warn"},
			},
			Choices: []Choice{
				{
					Label:    "[A] send it anyway — probably unmonitored at 3am.",
					Next:     "phish_success",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 5,
				},
				{
					Label:    "[B] find a different target — costs time, lower risk.",
					Next:     "legacy_probe",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 10,
				},
			},
		},

		"vpn_exploit": {
			Prompt:  "root@ghost:~/tools/exploit$",
			Mystery: "// project_ARGO — not in the brief. SIEM flag active.",
			Text: `CVE fires clean. Inside the VPN in four minutes.

Network map resolves: file shares, print server, a dev box
labelled project_ARGO.

Then a flag trips. "Anomalous login location" — their SIEM.
Automated. Unacknowledged. Probably ignored at 3am.

Probably.`,
			TargetEntries: []DossierEntry{
				{"project_ARGO", "Dev server. Restricted. Not mentioned in client brief.", "hi"},
			},
			SituationLines: []StatusLine{
				{"Inside VPN. Active session.", "ok"},
				{"SIEM flag: anomalous login — unacknowledged.", "warn"},
			},
			Choices: []Choice{
				{
					Label:    "[A] move fast — find the data before anyone checks.",
					Next:     "fast_search",
					Type:     ChoiceDanger,
					HeatCost: 2,
					TimeCost: 5,
				},
				{
					Label:    "[B] kill session, spoof local IP, reconnect clean.",
					Next:     "vpn_clean",
					Type:     ChoiceCaution,
					HeatCost: 0,
					TimeCost: 14,
				},
				{
					Label:    "[C] peek at project_ARGO before going for the target.",
					Next:     "argo_peek",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 8,
				},
			},
		},

		"legacy_probe": {
			Prompt:  "root@ghost:~/recon$",
			Mystery: "// employee 0001 — terminated 2019. same dept as the dev server.",
			Text: `Legacy subdomain loads slowly. Flash-based employee directory —
2009, somehow still live.

Unauthenticated API: GET /api/employees returns the full staff
list. 428 names, departments, extensions.

You scroll through. One entry stops you:

  ID: 0001 — NAME: [REDACTED]
  DEPT: ARGO RESEARCH — TERMINATED 2019`,
			TargetEntries: []DossierEntry{
				{"Legacy subdomain", "internal-legacy.meridian.net — unauthenticated API. Full staff exposed.", ""},
				{"ARGO Research", "Dept of one redacted employee — ID 0001, terminated 2019.", "hi"},
			},
			ContactEntries: []DossierEntry{
				{"Marcus Webb", "Junior sysadmin. Low-privilege. Active account confirmed.", ""},
			},
			SituationLines: []StatusLine{
				{"No alerts triggered.", "ok"},
				{"Employee list acquired. 428 entries.", "ok"},
			},
			Choices: []Choice{
				{
					Label:    "[A] use the employee list for a targeted phish.",
					Next:     "phish_success",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 10,
				},
				{
					Label:    "[B] dig deeper into the legacy API.",
					Next:     "legacy_deep",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 8,
				},
				{
					Label:    "[C] brute-force credentials using the list.",
					Next:     "brute_force",
					Type:     ChoiceDanger,
					HeatCost: 2,
					TimeCost: 4,
				},
			},
		},

		"legacy_deep": {
			Prompt:  "root@ghost:~/recon/legacy$",
			Mystery: "// \"the mirror\" — left deliberately. someone wanted this found.",
			Text: `Another endpoint: GET /api/internal/messages?id=0001

Internal messages from employee 0001. Sent days before termination
in 2019. Short thread. Last entry:

  "If anyone finds this: the ARGO data is real.
   They know it too. I've hidden a copy somewhere
   they won't find it easily.

   Look for the mirror."

Your hands are still for a moment.

This wasn't in the brief.`,
			TargetEntries: []DossierEntry{
				{"Employee 0001", "Redacted. ARGO Research. Terminated 2019. \"Look for the mirror.\"", "hi"},
				{"\"the mirror\"", "A hidden ARGO copy. Location unknown. Left intentionally.", "hi"},
			},
			SituationLines: []StatusLine{
				{"No alerts triggered.", "ok"},
				{"New variable: someone wanted this found.", "warn"},
			},
			Choices: []Choice{
				{
					Label:    "[A] ignore it — get the trial data and leave.",
					Next:     "win_clean",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 10,
				},
				{
					Label:    "[B] search for \"mirror\" across accessible file shares.",
					Next:     "find_mirror",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 16,
				},
				{
					Label:    "[C] save the message thread, then proceed.",
					Next:     "phish_success",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 5,
				},
			},
		},

		"phish_success": {
			Prompt:  "root@ghost:~$",
			Mystery: "// ARGO is ring-fenced. Marcus can see it. Can't open it.",
			Text: `Credentials land 18 minutes later. Marcus Webb — junior sysadmin.
Limited access, but inside.

Read permission on /projects/clinical/ — trial data. 4.2GB.
22 minutes to exfil at current speed. Window is getting narrow.

Next to it: /projects/ARGO — visible but inaccessible. Permissions
controlled separately from everything else.`,
			ContactEntries: []DossierEntry{
				{"Marcus Webb", "Creds acquired. Read: /projects/clinical. No ARGO access.", ""},
			},
			TargetEntries: []DossierEntry{
				{"Clinical trial data", "/projects/clinical — 4.2GB. 22min exfil at current speed.", ""},
				{"/projects/ARGO", "Visible. Restricted. Access-controlled separately.", "hi"},
			},
			SituationLines: []StatusLine{
				{"Active session as Marcus Webb.", "ok"},
				{"Exfil window: ~40 minutes remaining.", "warn"},
			},
			Choices: []Choice{
				{
					Label:    "[A] start the exfil — get what you came for.",
					Next:     "win_clean",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 22,
				},
				{
					Label:    "[B] try to escalate Marcus's privileges first.",
					Next:     "priv_esc",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 14,
				},
				{
					Label:    "[C] pull clinical data AND attempt ARGO. risky.",
					Next:     "argo_peek",
					Type:     ChoiceDanger,
					HeatCost: 2,
					TimeCost: 20,
				},
			},
		},

		"vpn_clean": {
			Prompt: "root@ghost:~$",
			Text: `14 minutes to spoof a local IP and reconnect. The SIEM alert
closes — resolved as a false positive.

Back inside, quieter now. File share is where you left it.`,
			SituationLines: []StatusLine{
				{"Session clean. SIEM alert resolved.", "ok"},
				{"Time lost on re-entry.", "warn"},
			},
			Choices: []Choice{
				{
					Label:    "[A] go straight for the clinical trial data.",
					Next:     "win_clean",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 20,
				},
				{
					Label:    "[B] check project_ARGO before pulling the file.",
					Next:     "argo_peek",
					Type:     ChoiceCaution,
					HeatCost: 1,
					TimeCost: 10,
				},
			},
		},

		"argo_peek": {
			Prompt:  "root@ghost:~/mnt/projects/ARGO$",
			Mystery: "// the job was never about the data. it was about making this disappear.",
			Text: `The directory loads.

It's not a product. It's a study — a suppressed Phase III trial.
Same drug. Different results. Dramatically different.

The clean version shows the drug works. This version shows it
doesn't. Three adverse events: two serious, one fatal. None
reported to the FDA.

Sign-offs at the bottom of every document:

  S. Okoye — Chief of Clinical Strategy.

Your client hired you to steal the data that buries what you're
reading right now.`,
			TargetEntries: []DossierEntry{
				{"ARGO — suppressed trial", "Phase III. Adverse events: 2 serious, 1 fatal. Not FDA-reported.", "danger"},
				{"Clinical trial data", "Positive results. The version your client wants.", "al"},
			},
			ContactEntries: []DossierEntry{
				{"S. Okoye", "Chief of Clinical Strategy. Signed every ARGO document.", "al"},
				{"Client", "knows this network better than a buyer should.", "al"},
			},
			SituationLines: []StatusLine{
				{"You know what this job actually is.", "warn"},
				{"Client wanted the cover-up completed.", "bad"},
			},
			Choices: []Choice{
				{
					Label:    "[A] complete the job — not your problem to fix.",
					Next:     "win_dirty",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 20,
				},
				{
					Label:    "[B] exfil the ARGO data instead. ghost your client.",
					Next:     "win_whistleblow",
					Type:     ChoiceSafe,
					HeatCost: 0,
					TimeCost: 20,
				},
				{
					Label:    "[C] take both datasets. figure it out later.",
					Next:     "win_both",
					Type:     ChoiceDanger,
					HeatCost: 1,
					TimeCost: 22,
				},
			},
		},

		// ── outcomes ────────────────────────────────────────────────────

		"find_mirror": {
			Prompt: "root@ghost:~/mnt/shares/archive$",
			TargetEntries: []DossierEntry{
				{"mirror_2019.tar.gz", "Hidden ARGO backup. Sign-offs trace to: Rexford Consulting.", "hi"},
				{"Rexford Consulting", "Third-party signatory. Not publicly linked to Meridian.", "hi"},
			},
			SituationLines: []StatusLine{
				{"Exfil complete. Session clean.", "ok"},
				{"You have more than you went in for.", "warn"},
			},
			Outcome: &Outcome{
				Type: "success",
				Text: `Search for "mirror" returns one hit:

  /archive/legacy_backup/mirror_2019.tar.gz

800MB. Encrypted — but the key is embedded in the filename.
Employee 0001 was thorough, not careful.

Inside: complete ARGO dataset, timestamped originals, with
internal sign-offs traced to a third party.

  Rexford Consulting.

Clean exit. 3 minutes before IT rotation.

Someone left this here on purpose. They wanted it found.
And now you have it — and a name you didn't go in with.`,
			},
		},

		"fast_search": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"SIEM flag escalated to human review.", "bad"},
				{"Session terminated remotely.", "bad"},
			},
			Outcome: &Outcome{
				Type: "fail",
				Text: `Folder found in three minutes. Exfil starts.

SIEM flag escalates. On-call engineer logs in remotely.
Active session — unrecognised location. He kills it.

60% complete. Incomplete data. No payment.

You left traces.`,
			},
		},

		"brute_force": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Account lockout triggered.", "bad"},
				{"IR team paged. You have hours.", "bad"},
			},
			Outcome: &Outcome{
				Type: "gameover",
				Text: `Account lockout after the 8th attempt.

Lockout generates a human-facing alert. Your VPN IP is
logged. Meridian's IR team is paged.

Game over.`,
			},
		},

		"priv_esc": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Escalation failed — logged to syslog.", "bad"},
				{"Exfil complete. Footprint left.", "warn"},
			},
			Outcome: &Outcome{
				Type: "fail",
				Text: `Escalation hits a patch from two weeks ago.
Failed attempt logs to syslog.

You pull the clinical data anyway. It completes.
But you left a footprint that surfaces in the morning audit.

Technically done. Not clean.`,
			},
		},

		"win_clean": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Exfil complete. No traces.", "ok"},
				{"You didn't look at what else was there.", "warn"},
			},
			Outcome: &Outcome{
				Type: "success",
				Text: `Exfil completes with 11 minutes to spare.
Session wiped. No alerts. Clean exit.

06:49 AM. You get paid.

But closing the laptop, you keep thinking about the
ARGO folder you didn't open.`,
			},
		},

		"win_dirty": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Job complete. Delivered.", "ok"},
				{"ARGO data remains buried.", "bad"},
			},
			Outcome: &Outcome{
				Type: "success",
				Text: `Clinical data pulled. Exfil complete. Silent exit.

Money arrives three days later.

You know what you handed over. You know what it'll
be used to suppress.

Some people sleep fine after that.`,
			},
		},

		"win_whistleblow": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Client ghosted. ARGO secured.", "ok"},
				{"Unknown contact made. Client now hostile.", "warn"},
			},
			Outcome: &Outcome{
				Type: "success",
				Text: `You pull the ARGO data and go silent. No delivery.
Client messages go unanswered.

Three weeks later an anonymous tip reaches the FDA.
Meridian stock drops 40%.

You didn't get paid. But the money was dirty.

Your client is looking for you. So is someone else —
unknown number, one message:

  "Thank you."`,
			},
		},

		"win_both": {
			Prompt: "root@ghost:~$",
			SituationLines: []StatusLine{
				{"Both files secured. Leverage acquired.", "ok"},
				{"Multiple parties will want these. You're exposed.", "bad"},
			},
			Outcome: &Outcome{
				Type: "success",
				Text: `Both datasets exfil with 6 minutes to spare.
Clean exit.

You now hold the data your client wants — and the
data that makes it worthless. Or priceless, depending
on who's buying.

This job just became something else entirely.

What you do next is the real game.`,
			},
		},
	}
}
