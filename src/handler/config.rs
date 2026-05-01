use std::fs;
use std::io::Write;
use std::path::PathBuf;
use config::Config;
use serde::Serialize;

fn config_file_path(file_name: String) -> PathBuf {
    let mut path = dirs::config_dir().unwrap();
    path.push("moose");
    path.push(file_name);
    path
}

pub fn ensure_config_exists<T>(file_name: String) -> Result<(), Box<dyn std::error::Error>>
where
    T: Serialize + Default,
{
    let path = config_file_path(file_name);

    if !path.exists() {
        if let Some(parent) = path.parent() {
            fs::create_dir_all(parent)?;
        }

        let default_config = T::default();
        let toml = toml::to_string_pretty(&default_config)?;

        let mut file = fs::File::create(path)?;
        file.write_all(toml.as_bytes())?;
    }

    Ok(())
}

pub fn make_config(file_name: String, required: bool) -> Result<Config, config::ConfigError> {
    let path = config_file_path(file_name);

    Config::builder()
        .add_source(config::File::from(path).required(required))
        .build()
}

pub fn vec_config(config: String) -> Vec<String> {
    let mut vec: Vec<String> = config.split('+').map(|s| s.trim().to_lowercase()).collect();
    vec.sort();
    vec
}

pub fn char_config(config: String) -> Option<char> {
    if config.starts_with("char:") {
        config.trim_start_matches("char:").chars().next()
    } else {
        None
    }
}