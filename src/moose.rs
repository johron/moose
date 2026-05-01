use crossterm::event::Event;
use serde::{Deserialize, Serialize};
use crate::handler::input::InputEvent;
use crate::handler::{input, workspace::Workspace};
use crate::handler::config::{char_config, ensure_config_exists, make_config};
use crate::handler::global_workspace::GlobalWorkspace;
use crate::panel::command_bar::command_bar::CommandBar;

#[derive(Debug)]
pub struct Moose {
    pub global_workspace: GlobalWorkspace,
    workspaces: Vec<Workspace>,
    active: usize,
    mode: MooseMode,
    should_quit: bool,
    pub config: MooseConfig,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct MooseConfig {
    pub enter_command_mode: String,
}

impl Default for MooseConfig {
    fn default() -> Self {
        MooseConfig {
            enter_command_mode: String::from("q"),
        }
    }
}

#[derive(Debug)]
enum MooseMode {
    Panel,
    GlobalPanel,
}

impl Moose {
    pub fn new() -> Self {
        Moose {
            global_workspace: GlobalWorkspace::new(),
            workspaces: Vec::new(),
            active: 0,
            mode: MooseMode::Panel,
            should_quit: false,
            config: MooseConfig::default(),
        }
    }

    pub fn init(&mut self) {
        let config = self.init_config();
        if config.is_ok() {
            self.config = config.unwrap();
        } else {
            eprintln!("Could not load moose config {:?}", config.err().unwrap());
        }

        self.global_workspace.set_bottom(Box::new(CommandBar::new()));
    }

    fn init_config(&self) -> Result<MooseConfig, config::ConfigError> {
        ensure_config_exists::<MooseConfig>(String::from("moose.toml")).expect("Failed to ensure config file exists");
        let cfg = make_config(String::from("moose.toml"), true)?;

        let config: MooseConfig = cfg.try_deserialize()?;
        Ok(config)
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
        self.should_quit
    }

    pub fn handle_terminal_event(&mut self, ev: Event) {
        if let Some(input) = input::from_crossterm_event(ev) {
            match self.mode {
                MooseMode::Panel => {
                    let config = self.config.clone();

                    if let Some(workspace) = self.active_workspace_mut() {
                        if let Some(panel) = workspace.active_panel_mut() {
                            if panel.is_normal_mode() {
                                match input {
                                    InputEvent::Char(c) => {
                                        if Some(c) == char_config(config.enter_command_mode) {
                                            todo!("Change to GlobalPanel mode, show command global panel")
                                        } else if Some(c) == char_config(String::from("char:-")) {
                                            self.should_quit = true;
                                        } else {
                                            panel.input(input);
                                        }
                                    },
                                    _ => panel.input(input)

                                }
                            } else {
                                panel.input(input);
                            }
                        }
                    }
                },
                MooseMode::GlobalPanel => {

                }
            }
        }
    }
}