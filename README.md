# GitHub Memory Bank

This repository contains a tool to help users implement a memory bank system for use with AI agents like GitHub Copilot. The system is designed to maintain context across sessions using specialized modes that handle different phases of the development process.

## Overview

The Memory Bank system utilizes a structured approach to manage tasks, context, and progress throughout a development workflow. It relies on a set of predefined "chat modes," each tailored for a specific stage:

*   **VAN Mode (Initialization):** Sets up the initial context, checks the memory bank status, and determines the task's complexity level.
*   **PM Mode (Product Requirements):** Guides the creation of detailed Product Requirements Documents (PRDs) through structured clarifying questions and comprehensive documentation.
*   **PLAN Mode (Task Planning):** Creates a detailed plan for task execution based on the complexity.
*   **CREATIVE Mode (Design Decisions):** Facilitates detailed design and architecture work for components flagged during planning.
*   **IMPLEMENT Mode (Code Implementation):** Guides the building of planned changes, following the implementation plan and creative phase decisions.
*   **REFLECT+ARCHIVE Mode (Review & Archiving):** Facilitates reflection on the completed task and then archives relevant documentation, updating the Memory Bank.

## Key Features

*   **Mode-Specific Workflows:** Each mode has a defined process, often visualized with Mermaid diagrams, to guide the user and the AI.
*   **Contextual Memory:** The system uses a "Memory Bank" (a set of structured markdown files) to maintain context, track tasks (`tasks.md`), manage active context (`activeContext.md`), store creative decisions (`creative-*.md`), and log progress (`progress.md`).
*   **Isolation-Focused Design:** Each mode is designed to load only the rules and information it needs, optimizing for efficiency.
*   **Structured Documentation:** The system emphasizes clear documentation at each stage, including planning documents, design guidelines, implementation details, and reflection notes.
*   **Verification Checkpoints:** Each mode includes verification steps to ensure that processes are followed correctly and that the Memory Bank is updated appropriately.

## How It Works

The system is intended to be used with GitHub Copilot Chat's custom modes. The user interacts with the AI, and the AI follows the instructions and workflows defined in the `.chatmode.md` files.

1.  **Initialization (VAN Mode):** A new task begins in VAN mode, where the project brief and task complexity are established.
2.  **Product Requirements (PM Mode):** When needed, this mode creates comprehensive Product Requirements Documents (PRDs) through structured questioning, ensuring clear feature definition before planning begins.
3.  **Planning (PLAN Mode):** Based on the task's complexity, a detailed implementation plan is created. Components requiring significant design work are flagged for a "Creative Phase."
4.  **Creative Design (CREATIVE Mode):** If needed, this mode is used to explore design options (architecture, algorithms, UI/UX) for flagged components.
5.  **Implementation (IMPLEMENT Mode):** Code changes are made according to the plan and any creative design decisions.
6.  **Reflection and Archiving (REFLECT+ARCHIVE Mode):** After implementation, the work is reviewed, lessons are learned, and all relevant documentation is archived. The Memory Bank is updated to reflect the completed task.

## Getting Started

*(This section would typically include instructions on how to set up and use the tool. This can be filled in as the tool is developed.)*

## Contributing

*(This section would typically include guidelines for contributing to the project. This can be filled in as the tool is developed.)*

---

This README is based on the configuration files found in the `.github` directory, including:
*   `.github/chatmodes/⚒️ IMPLEMENT.chatmode.md`
*   `.github/chatmodes/🎨 CREATIVE.chatmode.md`
*   `.github/chatmodes/📋 PLAN.chatmode.md`
*   `.github/chatmodes/� PM.chatmode.md`
*   `.github/chatmodes/�🔍 REFLECT.chatmode.md`
*   `.github/chatmodes/🔍 VAN.chatmode.md`
*   `.github/instructions/main.instructions.md`
