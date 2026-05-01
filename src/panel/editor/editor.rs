use crate::handler::input::InputEvent;
use crate::panel::editor::config::{init_config, EditorConfig};
use crate::panel::editor::event::handle_event;
use crate::panel::editor::input::handle_input;
use crate::panel::editor::renderer::render;
use crate::panel::panel::{Cursor, Panel};
use ratatui::layout::Rect;
use ratatui::Frame;
use ropey::Rope;

#[derive(Debug)]
pub struct Editor {
    init: bool,
    pub rope: Rope,
    pub scroll_offset: usize,
    pub cursors: Vec<Cursor>,
    pub mode: EditorMode,
    pub config: EditorConfig,
}

#[derive(Debug)]
pub enum EditorMode {
    Normal,
    Insert,
    // ++
}

impl Panel for Editor {
    fn init(&mut self) {
        if self.init {
            return;
        }
        let config = init_config();
        if config.is_ok() {
            self.config = config.unwrap();
        } else {
            eprintln!("Could not load editor config {:?}", config.err().unwrap());
        }

        self.init = true;
    }

    fn is_initialized(&self) -> bool {
        self.init
    }
    fn is_normal_mode(&self) -> bool {
        matches!(self.mode, EditorMode::Normal)
    }

    fn identity(&self) -> &str {
        "builtin:editor"
    }

    fn title(&self) -> String {
        String::from("Editor")
    }

    fn render(&self, frame: &mut Frame, area: Rect) {
        render(self, frame, area);
    }

    fn input(&mut self, input: InputEvent) {
        let event = handle_input(self, input);
        handle_event(self, event);
    }
}

impl Editor {
    pub fn new() -> Self {
        Editor {
            init: false,
            rope: Rope::new(),
            scroll_offset: 0,
            cursors: vec![Cursor::new()],
            mode: EditorMode::Normal,
            config: EditorConfig::default(),
        }
    }
}