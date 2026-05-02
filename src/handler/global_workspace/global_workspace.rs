use crate::handler::global_workspace::renderer::render;
use crate::handler::workspace::workspace::Workspace;
use crate::panel::global_panel::GlobalPanel;
use ratatui::layout::Rect;
use ratatui::Frame;
use crate::handler::global_workspace::config::{init_config, MooseConfig};

#[derive(Debug)]
pub struct GlobalWorkspace {
    pub bottom: Option<Box<dyn GlobalPanel>>,
    pub left: Option<Box<dyn GlobalPanel>>,
    pub right: Option<Box<dyn GlobalPanel>>,
    pub active: GlobalWorkspaceActive,
    pub config: MooseConfig,
    pub should_quit: bool,
}

#[derive(Debug)]
pub enum GlobalWorkspaceActive {
    Workspace,
    Bottom,
    Left,
    Right,
}

impl GlobalWorkspace {
    pub fn new() -> Self {
        GlobalWorkspace {
            bottom: None,
            left: None,
            right: None,
            active: GlobalWorkspaceActive::Workspace,
            config: MooseConfig::default(),
            should_quit: false,
        }
    }

    pub fn init(&mut self) {
        let config = init_config();
        if config.is_ok() {
            self.config = config.unwrap();
        } else {
            eprintln!("Could not load moose config {:?}", config.err().unwrap());
        }
    }
    
    pub fn set_bottom(&mut self, mut panel: Box<dyn GlobalPanel>) {
        panel.init();
        self.bottom = Some(panel);
    }

    pub fn set_left(&mut self, mut panel: Box<dyn GlobalPanel>) {
        panel.init();
        self.left = Some(panel);
    }

    pub fn set_right(&mut self, mut panel: Box<dyn GlobalPanel>) {
        panel.init();
        self.right = Some(panel);
    }

    pub fn get_active_panel(&self) -> Option<&Box<dyn GlobalPanel>> {
        match self.active {
            GlobalWorkspaceActive::Bottom => if let Some(bottom) = &self.bottom {
                Some(&bottom)
            } else {
                None
            },
            GlobalWorkspaceActive::Left => if let Some(left) = &self.left {
                Some(&left)
            } else {
                None
            },
            GlobalWorkspaceActive::Right => if let Some(right) = &self.right {
                Some(&right)
            } else {
                None
            },
            _ => None,
        }
    }

    pub fn render(&self, child_workspace: Option<&Workspace>, frame: &mut Frame, area: Rect) {
        render(self, child_workspace, frame, area);
    }
    
    pub fn tick(&mut self) {
        if let Some(bottom) = self.bottom.as_mut() {
            if bottom.should_quit() {
                self.bottom = None;
                self.active = GlobalWorkspaceActive::Workspace;
            }
        } else if let Some(left) = self.left.as_mut() {
            if left.should_quit() {
                self.left = None;
                self.active = GlobalWorkspaceActive::Workspace;
            }
        } else if let Some(right) = self.right.as_mut() {
            if right.should_quit() {
                self.right = None;
                self.active = GlobalWorkspaceActive::Workspace;
            }
        }
    }
}