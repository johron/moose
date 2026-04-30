use crate::handler::input::InputEvent;
use crate::panel::editor::event::handle_event;
use crate::panel::editor::input::handle_input;
use crate::panel::editor::renderer::render;
use crate::panel::panel::Panel;
use ratatui::layout::Rect;
use ratatui::Frame;
use ropey::Rope;
use crate::panel::editor::config::{init_config, EditorConfig};

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
#[derive(Clone)]
pub struct Cursor {
    pub line: usize,
    pub col: usize,
    pub goal: usize,
}

impl Cursor {
    pub fn new() -> Self {
        Cursor {
            line: 0,
            col: 0,
            goal: 0,
        }
    }

    pub fn from(line: usize, col: usize, goal: usize) -> Self {
        Cursor {
            line,
            col,
            goal,
        }
    }
}

#[derive(Debug)]
pub enum EditorMode {
    Normal,
    Insert,
    Command,
    // ++
}

impl Panel for Editor {
    fn init(&mut self) {
        if self.init {
            return;
        }
        println!("Created a new editor panel {}", self.title());
        let config = init_config();
        if config.is_ok() {
            self.config = config.unwrap();
        } else {
            eprintln!("Could not load config.. {:?}", config.err().unwrap());
        }

        self.init = true;
    }

    fn is_init(&self) -> bool {
        self.init
    }

    fn identity(&self) -> &str {
        "editor"
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

#[derive(Debug)]
pub enum EditorEdit {
    Insert(char),
    Delete,
    Newline,
    Tab,
    Backtab
}

#[derive(Debug)]
pub enum EditorEvent {
    Edit(EditorEdit), // relative
    CursorMove(usize, usize), // x, y, relative
    Scroll(usize, usize), // x, y, relative
    Mode(EditorMode),
    Quit,
    None,
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