use std::fmt::{Debug, Formatter};
use crate::handler::input::InputEvent;
use crate::panel::command_bar::config::{init_config, CommandBarConfig};
use crate::panel::command_bar::renderer::render;
use crate::panel::global_panel::{GlobalPanel};
use crate::panel::panel::{Cursor, Panel};
use ratatui::layout::Rect;
use ratatui::Frame;
use crate::handler::config::{char_config, vec_config};

#[derive(Debug)]
pub struct CommandBar {
    init: bool,
    pub cmd: String,
    pub cursor: Cursor,
    pub config: CommandBarConfig,
    pub mode: CommandBarMode,
    show: bool,
    floating: bool,
    pub should_quit: bool,
}

#[derive(Debug, Eq, Hash, PartialEq)]
pub enum CommandBarMode {
    Normal,
    Insert,
}

impl CommandBar {
    pub fn new() -> Self {
        CommandBar {
            init: false,
            cmd: String::new(),
            cursor: Cursor::new(),
            config: CommandBarConfig::default(),
            mode: CommandBarMode::Insert,
            show: false,
            floating: false,
            should_quit: false,
        }
    }

    pub fn parse_command(&mut self) -> Option<ParsedCommand> {
        let mut parts = self.cmd.split_whitespace();
        let name = parts.next()?.to_string();
        let args = parts.map(|s| s.to_string()).collect();
        Some(ParsedCommand { name, args })
    }
}

pub struct ParsedCommand {
    pub name: String,
    pub args: Vec<String>,
}

impl Panel for CommandBar {
    fn init(&mut self) {
        let config = init_config();
        if config.is_ok() {
            self.config = config.unwrap();

            if self.config.start_in_insert_mode {
                self.mode = CommandBarMode::Insert;
            } else {
                self.mode = CommandBarMode::Normal;
            }
        } else {
            eprintln!("Could not load command bar config {:?}", config.err().unwrap());
        }
    }

    fn is_initialized(&self) -> bool {
        self.init
    }

    fn is_normal_mode(&self) -> bool {
        true
    }

    fn identity(&self) -> &str {
        "builtin:command_bar"
    }

    fn title(&self) -> String {
        String::from("Command Bar")
    }

    fn render(&self, frame: &mut Frame, area: Rect) {
        render(self, frame, area);
    }

    fn input(&mut self, input: InputEvent) {
        match self.mode {
            CommandBarMode::Normal => {
                match input {
                    InputEvent::Char(c) => {
                        if char_config(self.config.enter_insert_mode.clone()) == Some(c) {
                            self.mode = CommandBarMode::Insert;
                        } else if char_config(self.config.quit_command_bar.clone()) == Some(c) {
                            self.should_quit = true;
                        }
                    },
                    InputEvent::Keyboard(vec) => {
                        if vec_config(self.config.enter_insert_mode.clone()) == vec {
                            self.mode = CommandBarMode::Insert;
                        } else if vec_config(self.config.quit_command_bar.clone()) == vec {
                            self.should_quit = true;
                        }
                    },
                    _ => {}
                }
            },
            CommandBarMode::Insert => {
                match input {
                    InputEvent::Keyboard(vec) => {
                        if vec_config(self.config.execute_command.clone()) == vec {
                            todo!()
                            // if self.config.enter_normal_mode_after_execute ..
                        } else if vec_config(self.config.delete_char.clone()) == vec {
                            // todo: respect cursor
                            self.cmd.pop();
                        } else if vec_config(self.config.enter_normal_mode.clone()) == vec {
                            self.mode = CommandBarMode::Normal;
                        } else {}
                    },
                    InputEvent::Char(c) => {
                        // make it better, align with cursor
                        self.cmd.push(c)
                    },
                    _ => {}
                }
            }
        }
    }

    fn should_quit(&self) -> bool {
        self.should_quit
    }
}

impl GlobalPanel for CommandBar {
    fn set_floating(&mut self, floating: bool) {
        self.floating = floating;
    }

    fn is_floating(&self) -> bool {
        self.floating
    }
}
