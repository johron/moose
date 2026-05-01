# Moose
```
                ___            ___
               /   /          /   \        
               \_   \      \ /  __/
                _\   \      \  /__          
                \___  \____/   __/         
                    \_       _/            Modern mode-oriented configurable and
                      |0   0 \_            extendable terminal-based text editor
                      |                    
                    _/    /\
                   /o)  (o/\ \_
                   \_____/ /
                     \_-__/ 
```

## TODO
- [ ] Finish editor features, newlines, tabs, cursor movement, scroll, modes
- [ ] More advanced rendering
  - [ ] Drawing multiple panels in one workspace, layout etc
  - [ ] Global panels: Layout things. Bottom takes all width. Side panels take height - height of bottom panel. Main panel takes rest of space.
  - [ ] kommandoregister der hvert panel kan registrere commands med en callback. Legg til funksjon og Panel trait: commands som returnerer en Hashmap(cmd: str, callback: func)
- [ ] Lua plugin system
  - [ ] To start: Syntax highlighting for editor panel. Editor Plugins. Editor panel needs to expose a plugin API?
  - [ ] Later: Custom panels, shortcuts, .., 