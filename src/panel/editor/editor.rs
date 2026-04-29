use ratatui::Frame;
use ratatui::layout::Rect;
use ropey::Rope;
use crate::handler::input::InputEvent;
use crate::panel::editor::renderer::render;
use crate::panel::panel::Panel;

#[derive(Debug)]
pub struct Editor {
    init: bool,
    rope: Rope,
    scroll_offset: usize,
    cursors: Vec<Cursor>,
    mode: EditorMode,
}

#[derive(Debug)]
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
enum EditorMode {
    Normal,
    Insert,
    Command,
    // ++
}

impl Panel for Editor {
    fn init(&mut self) {
        if !self.init {
            println!("Created a new editor panel {}", self.title());
            self.init = true;
        }
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
        todo!()
    }
}

impl Editor {
    pub fn new() -> Self {
        Editor {
            init: false,
            rope: Rope::new(),
            scroll_offset: 0,
            cursors: Vec::new(),
            mode: EditorMode::Normal,
        }
    }

    pub fn get_rope(&self) -> &Rope {
        &self.rope
    }

    pub fn get_scroll_offset(&self) -> usize {
        self.scroll_offset
    }

    pub fn get_cursors(&self) -> &Vec<Cursor> {
        &self.cursors
    }
}