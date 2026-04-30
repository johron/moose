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
- [x] Input
- [x] Config system
- [ ] Finish editor features, newlines, tabs, cursor movement, scroll, modes
- [ ] Figure out global modes? how will this work??
- [ ] More advanced rendering
  - [ ] Drawing multiple panels in one workspace, layout etc
  - [ ] Command bar, should this be global, or managed per panel? registered commands may only be per panel.
    - [ ] `shift + q` for å åpne global command bar. Eller så har jeg et kommandoregister der hvert panel kan registrere commands med en callback.
- [ ] Lua plugin system
  - [ ] To start: Syntax highlighting for editor panel. Editor Plugins. Editor panel needs to expose a plugin API?
  - [ ] Later: Custom panels, shortcuts, .., 