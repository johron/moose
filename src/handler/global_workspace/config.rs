use serde::{Deserialize, Serialize};
use crate::handler::config::{ensure_config_exists, make_config};

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct MooseConfig {
    pub enter_global_command_mode: String,
}

impl Default for MooseConfig {
    fn default() -> Self {
        MooseConfig {
            enter_global_command_mode: String::from("q"),
        }
    }
}

pub fn init_config() -> Result<MooseConfig, config::ConfigError> {
    ensure_config_exists::<MooseConfig>(String::from("moose.toml")).expect("Failed to ensure config file exists");
    let cfg = make_config(String::from("moose.toml"), true)?;

    let config: MooseConfig = cfg.try_deserialize()?;
    Ok(config)
}