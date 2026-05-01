use crossterm::event::Event;
use serde::{Deserialize, Serialize};
use crate::handler::input::InputEvent;
use crate::handler::{input, workspace::Workspace};
use crate::handler::config::{char_config, ensure_config_exists, make_config};

#[derive(Debug)]
pub struct Moose {
    workspaces: Vec<Workspace>,
    active: usize,
    mode: GlobalMode,
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
enum GlobalMode {
    Panel,
    GlobalPanel,
}

impl Moose {
    pub fn new() -> Self {
        Moose {
            workspaces: Vec::new(),
            active: 0,
            mode: GlobalMode::Panel,
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

    pub fn active_workspace(&mut self) -> Option<&mut Workspace> {
        self.workspaces.get_mut(self.active)
    }

    pub fn should_quit(&self) -> bool {
        self.should_quit
    }

    pub fn handle_terminal_event(&mut self, ev: Event) {
        if let Some(input) = input::from_crossterm_event(ev) {
            match self.mode {
                GlobalMode::Panel => {
                    let config = self.config.clone();

                    if let Some(workspace) = self.active_workspace() {
                        if let Some(panel) = workspace.active_panel() {
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
                GlobalMode::GlobalPanel => {

                }
            }
        }
    }
}