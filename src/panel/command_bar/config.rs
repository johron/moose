use crate::handler::config::{ensure_config_exists, make_config};
use serde::{Deserialize, Serialize};

pub fn init_config() -> Result<CommandBarConfig, config::ConfigError> {
    ensure_config_exists::<CommandBarConfig>(String::from("builtin/command_bar.toml")).expect("Failed to ensure config file exists");
    let cfg = make_config(String::from("builtin/command_bar.toml"), true)?;

    let config: CommandBarConfig = cfg.try_deserialize()?;
    Ok(config)
}

#[derive(Debug, Deserialize, Serialize, Clone, Eq, Hash, PartialEq)]
pub struct CommandBarConfig {
    pub start_in_insert_mode: bool,
    pub enter_normal_mode_after_exec: bool,
    
    pub quit_command_bar: String,
    pub enter_insert_mode: String,
    pub enter_normal_mode: String,
    
    pub execute_command: String,
    pub delete_char: String,
    pub cursor_left: String,
    pub cursor_right: String,
}

impl Default for CommandBarConfig {
    fn default() -> CommandBarConfig {
        CommandBarConfig {
            start_in_insert_mode: true,
            enter_normal_mode_after_exec: true,
            
            quit_command_bar: String::from("char:q"),
            enter_insert_mode: String::from("char:i"),
            enter_normal_mode: String::from("esc"),
            
            execute_command: String::from("enter"),
            delete_char: String::from("backspace"),
            cursor_left: String::from("left"),
            cursor_right: String::from("right"),
        }
    }
}