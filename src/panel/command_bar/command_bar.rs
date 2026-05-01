use ratatui::Frame;
use ratatui::layout::Rect;
use serde::{Deserialize, Serialize};
use crate::panel::global_panel::GlobalPanel;
use crate::handler::input::InputEvent;
use crate::panel::command_bar::renderer::render;
use crate::panel::panel::{Cursor, Panel};

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
        todo!()
    }

    fn is_initialized(&self) -> bool {
        self.init
    }

    fn is_normal_mode(&self) -> bool {
        true
    }

    fn identity(&self) -> &str {
        "moose:command_bar"
    }

    fn title(&self) -> String {
        String::from("Command Bar")
    }

    fn render(&self, frame: &mut Frame, area: Rect) {
        render(self, frame, area);
    }

    fn input(&mut self, input: InputEvent) {
        todo!()
    }
}

impl GlobalPanel for CommandBar {
    fn set_show(&mut self, show: bool) {
        self.show = show;
    }

    fn is_shown(&self) -> bool {
        self.show
    }

    fn set_floating(&mut self, floating: bool) {
        self.floating = floating;
    }

    fn is_floating(&self) -> bool {
        self.floating
    }
}

#[derive(Debug, Deserialize, Serialize, Clone, Eq, Hash, PartialEq)]
pub struct CommandBarConfig {

}

impl Default for CommandBarConfig {
    fn default() -> CommandBarConfig {
        CommandBarConfig {

        }
    }
}