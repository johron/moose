use std::fmt::Debug;
use crate::panel::panel::Panel;

pub trait GlobalPanel: Debug + Panel {
    fn set_show(&mut self, show: bool);
    fn is_shown(&self) -> bool;

    fn set_floating(&mut self, floating: bool);
    fn is_floating(&self) -> bool;
}