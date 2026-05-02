use std::any::Any;
use std::collections::HashMap;
use crate::panel::panel::Panel;
use std::fmt::Debug;
use crate::handler::global_workspace::global_workspace::GlobalWorkspaceActive;

pub trait GlobalPanel: Debug + Panel + Any {
    //fn set_show(&mut self, show: bool);
    //fn is_shown(&self) -> bool;
    fn set_floating(&mut self, floating: bool);
    fn is_floating(&self) -> bool;
}

pub struct GlobalPanelMeta {
    pub panel: Box<dyn GlobalPanel>,
    pub toggle_show_shortcut: String,
    pub location: GlobalWorkspaceActive,
}