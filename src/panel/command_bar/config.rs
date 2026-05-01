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
    pub execute_command: String,
    pub cursor_left: String,
    pub cursor_right: String,
}

impl Default for CommandBarConfig {
    fn default() -> CommandBarConfig {
        CommandBarConfig {
            execute_command: String::from("enter"),
            cursor_left: String::from("left"),
            cursor_right: String::from("right"),
        }
    }
}