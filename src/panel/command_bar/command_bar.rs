use crate::handler::input::InputEvent;
use crate::panel::command_bar::config::{init_config, CommandBarConfig};
use crate::panel::command_bar::renderer::render;
use crate::panel::global_panel::{GlobalPanel, GlobalPanelMeta};
use crate::panel::panel::{Cursor, Panel};
use ratatui::layout::Rect;
use ratatui::Frame;
use crate::handler::config::vec_config;

#[derive(Debug, Eq, Hash, PartialEq)]
pub struct CommandBar {
    init: bool,
    pub cmd: String,
    pub cursor: Cursor,
    pub config: CommandBarConfig,
    show: bool,
    floating: bool,
}

impl CommandBar {
    pub fn new() -> Self {
        CommandBar {
            init: false,
            cmd: String::new(),
            cursor: Cursor::new(),
            config: CommandBarConfig::default(),
            show: false,
            floating: false,
        }
    }
}

impl Panel for CommandBar {
    fn init(&mut self) {
        let config = init_config();
        if config.is_ok() {
            self.config = config.unwrap();
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
        match input {
            InputEvent::Keyboard(vec) => {
                if vec_config(self.config.execute_command.clone()) == vec {
                    todo!()
                } else if vec_config(self.config.delete_char.clone()) == vec {
                    self.cmd.pop();
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

impl GlobalPanel for CommandBar {
    //fn set_show(&mut self, show: bool) {
    //    self.show = show;
    //}

    //fn is_shown(&self) -> bool {
    //    self.show
    //}

    fn set_floating(&mut self, floating: bool) {
        self.floating = floating;
    }

    fn is_floating(&self) -> bool {
        self.floating
    }
}