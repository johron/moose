use std::collections::HashMap;
use crate::panel::panel::Panel;
use std::fmt::Debug;

pub trait GlobalPanel: Debug + Panel {
    fn set_show(&mut self, show: bool);
    fn is_shown(&self) -> bool;

    fn set_floating(&mut self, floating: bool);
    fn is_floating(&self) -> bool;
    fn set_commands_hashmap(&mut self, commands_hashmap: HashMap<String, fn(>);
}