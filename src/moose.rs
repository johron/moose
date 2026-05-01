use crate::handler::global_workspace::global_workspace::GlobalWorkspace;
use crate::handler::{input, workspace::workspace::Workspace};
use crate::panel::command_bar::command_bar::CommandBar;
use crossterm::event::Event;

#[derive(Debug)]
pub struct Moose {
    pub global_workspace: GlobalWorkspace,
    workspaces: Vec<Workspace>,
    active: usize,
}

impl Moose {
    pub fn new() -> Self {
        Moose {
            global_workspace: GlobalWorkspace::new(),
            workspaces: Vec::new(),
            active: 0,
        }
    }

    pub fn init(&mut self) {
        self.global_workspace.init();
        self.global_workspace.set_bottom(Box::new(CommandBar::new()));
    }

    pub fn add_workspace(&mut self, make_active: bool) {
        self.workspaces.push(Workspace::new());
        if make_active {
            self.active = self.workspaces.len() - 1;
        }
    }

    pub fn active_workspace_mut(&mut self) -> Option<&mut Workspace> {
        self.workspaces.get_mut(self.active)
    }

    pub fn active_workspace(&self) -> Option<&Workspace> {
        self.workspaces.get(self.active)
    }

    pub fn should_quit(&self) -> bool {
        self.global_workspace.should_quit
    }

    pub fn handle_terminal_event(&mut self, ev: Event) {
        if let Some(input_event) = input::from_crossterm_event(ev) {
            let Moose {
                global_workspace,
                workspaces,
                active,
            } = self;

            let active_ws = workspaces.get_mut(*active);
            global_workspace.input(active_ws, input_event);
        }
    }
}