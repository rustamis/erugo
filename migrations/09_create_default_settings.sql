INSERT INTO
  settings (
    id,
    value,
    setting_group,
    created_at,
    updated_at
  )
VALUES
  (
    'max_file_size',
    '2G',
    'general',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'max_file_count',
    '50',
    'general',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'share_expiration_days',
    '7',
    'general',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'css_primary_color',
    'rgb(238, 193, 84)',
    'ui.css',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'css_secondary_color',
    'rgb(34, 34, 34)',
    'ui.css',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'css_accent_color',
    'rgb(84, 129, 238)',
    'ui.css',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'css_accent_color_light',
    'rgb(238, 238, 238)',
    'ui.css',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'logo_data',
    '',
    'ui',
    DATETIME('now'),
    DATETIME('now')
  ),
  (
    'application_name',
    'erugo File Sharing',
    'general',
    DATETIME('now'),
    DATETIME('now')
  );