<!doctype html>
<html lang="en">

<head>
  <meta charset="UTF-8" />
  <link rel="icon" type="image/svg+xml" href="/vite.svg" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap/5.0.2/css/bootstrap.min.css" />
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" integrity="sha512-Evv84Mr4kqVGRNSgIGL/F/aIDqQb7xQ2vcrdIwxfjThSH8CSR7PBEakCr51Ck+w+/U6swU2Im1vVX0SVk9ABhg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
  <title>{{ $settings['application_name'] }}</title>
  @vite('resources/js/app.js')
 
</head>

<body data-settings='{!! json_encode($settings) !!}'>
<style id="erugo-css-variables">
    :root {
      --primary-color: {{ $settings['css_primary_color'] }};
      --secondary-color: {{ $settings['css_secondary_color'] }};
      --accent-color: {{ $settings['css_accent_color'] }};
      --accent-color-light: {{ $settings['css_accent_color_light'] }};
    }
  </style>
  <div id="app"></div>
</body>

</html>