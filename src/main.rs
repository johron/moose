use crate::moose::Moose;
use crate::panel::editor::editor::Editor;
use crossterm::event::{DisableMouseCapture, KeyboardEnhancementFlags, PopKeyboardEnhancementFlags, PushKeyboardEnhancementFlags};
use crossterm::terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen};
use ratatui::backend::CrosstermBackend;
use ratatui::Terminal;
use std::io;
use std::io::StdoutLock;
use std::time::Duration;

mod panel;
mod moose;
mod handler;

fn main() -> io::Result<()> {
    let mut moose = Moose::new();
    moose.init();
    moose.add_workspace(true);

    if let Some(workspace) = moose.active_workspace() {
        workspace.add_panel(Box::new(Editor::new()), true);
    }

    let stdout = io::stdout();
    let mut stdout = stdout.lock();

    enable_raw_mode()?;

    #[cfg(not(windows))] // keyboard enhancements don't work on windows
    crossterm::execute!(stdout, EnterAlternateScreen, /*EnableMouseCapture remove for now, */PushKeyboardEnhancementFlags(
        KeyboardEnhancementFlags::DISAMBIGUATE_ESCAPE_CODES
    ))?;

    #[cfg(windows)]
    crossterm::execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;

    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    let res = run(&mut terminal, moose);

    disable_raw_mode()?;

    #[cfg(not(windows))]
    crossterm::execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture,
        PopKeyboardEnhancementFlags // dette ser ikke ut som at det fungerer.. kan ikke lukke nano etter å ha kjørt mos
    )?;

    #[cfg(windows)]
    crossterm::execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture
    )?;

    terminal.show_cursor()?;

    res
}

fn run(terminal: &mut Terminal<CrosstermBackend<StdoutLock>>, mut moose: Moose) -> io::Result<()> {
    loop {
        if moose.should_quit() {
            break Ok(());
        }

        while crossterm::event::poll(Duration::from_millis(10))? {
            let ev = crossterm::event::read()?;
            moose.handle_terminal_event(ev);
        }

        terminal.draw(|frame| {
            let area = frame.area();
            if let Some(workspace) = moose.active_workspace() {
                workspace.render(frame, area);
            }
        })?;
    }
}