#set-option -sa terminal-overrides ",xterm*:Tc"
set -g default-terminal "screen-256color"
set -g mouse on
set-environment -g PATH "/opt/homebrew/bin:/bin:/usr/bin"
unbind C-b
set -g prefix C-a
bind C-a send-prefix

# Vim style pane selection
bind h select-pane -L
bind j select-pane -D 
bind k select-pane -U
bind l select-pane -R

# Start windows and panes at 1, not 0
set -g base-index 1
set -g pane-base-index 1
set-window-option -g pane-base-index 1
set-option -g renumber-windows on

# Use Alt-arrow keys without prefix key to switch panes
bind -n M-Left select-pane -L
bind -n M-Right select-pane -R
bind -n M-Up select-pane -U
bind -n M-Down select-pane -D

# Shift arrow to switch windows
bind -n S-Left  previous-window
bind -n S-Right next-window

# Shift Alt vim keys to switch windows
bind -n M-H previous-window
bind -n M-L next-window
bind-key & kill-window
bind-key x kill-pane

set -g history-limit 9999
set-option -g default-shell '/bin/zsh'
set-option -g display-time 10000
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'christoomey/vim-tmux-navigator'
set -g @plugin 'tmux-plugins/tmux-yank'
#set -g @plugin 'seebi/tmux-colors-solarized'
set -g @plugin 'mkoga/tmux-solarized'
#set -g @plugin 'jimeh/tmux-themepack'
#set -g @plugin 'tmux-plugins/tmux-resurrect' # persist tmux sessions after computer restart
#set -g @plugin 'tmux-plugins/tmux-continuum' # automatically saves sessions for you every 15 minutes
set -g @plugin 'omerxx/tmux-sessionx'
#set -g @themepack 'powerline/double/cyan'
#set -g @plugin 'wfxr/tmux-power'
#set -g @plugin 'catppuccin/tmux'
#set -g @catppuccin_window_default_text "#W" # use "#W" for application instead of directory
#set -g @resurrect-capture-pane-contents 'on'
#set -g @continuum-restore 'on'
#set -g @continuum-save-interval '60'
#set -g @continuum-save-uptime 'on'
#set -g @colors-solarized 'dark'

set -g @sessionx-bind 'z'
run '~/.tmux/plugins/tpm/tpm'
# I recommend using `o` if not already in use, for least key strokes when launching
# set vi-mode
set-window-option -g mode-keys vi
# keybindings
bind-key -T copy-mode-vi v send-keys -X begin-selection
bind-key -T copy-mode-vi C-v send-keys -X rectangle-toggle
bind-key -T copy-mode-vi y send-keys -X copy-selection-and-cancel
bind-key g set-window-option synchronize-panes\; display-message "synchronize-panes is now #{?pane_synchronized,on,off}"
bind '"' split-window -v -c "#{pane_current_path}"
bind % split-window -h -c "#{pane_current_path}"
