use crossterm::event::Event;

use crate::handler::{input, workspace::Workspace};

#[derive(Debug)]
pub struct Moose {
    workspaces: Vec<Workspace>,
    active: usize,
    mode: GlobalMode,
    should_quit: bool,
}

#[derive(Debug)]
enum GlobalMode {
    Normal,
    Panel,
}

impl Moose {
    pub fn new() -> Self {
        Moose {
            workspaces: Vec::new(),
            active: 0,
            mode: GlobalMode::Panel,
            should_quit: false,
        }
    }

    pub fn add_workspace(&mut self, make_active: bool) {
        self.workspaces.push(Workspace::new());
        if make_active {
            self.active = self.workspaces.len() - 1;
        }
    }

    pub fn active_workspace(&mut self) -> Option<&mut Workspace> {
        self.workspaces.get_mut(self.active)
    }

    pub fn should_quit(&self) -> bool {
        self.should_quit
    }

    pub fn handle_terminal_event(&mut self, ev: Event) {
        println!("{:?}", ev);
        if let Some(input) = input::from_crossterm_event(ev) {
            match self.mode {
                GlobalMode::Panel => {
                    if let Some(workspace) = self.active_workspace() {
                        if let Some(panel) = workspace.active_panel() {
                            panel.input(input);
                        }
                    }
                },
                GlobalMode::Normal => {} 
            }
        }
    }
}