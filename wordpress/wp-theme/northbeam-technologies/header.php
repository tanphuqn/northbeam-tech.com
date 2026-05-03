<!doctype html>
<html <?php language_attributes(); ?>>

<head>
  <meta charset="<?php bloginfo('charset'); ?>" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="icon" type="image/svg+xml" href="<?php echo esc_url(get_template_directory_uri() . '/assets/favicon.svg'); ?>">
  <?php wp_head(); ?>
</head>

<body <?php body_class(); ?>>
  <?php wp_body_open(); ?>
  <header>
    <div class="container nav-wrap">
      <a class="logo" href="<?php echo esc_url(home_url('/')); ?>">
        <img src="<?php echo esc_url(get_template_directory_uri() . '/assets/logo.svg'); ?>" alt="Northbeam Technologies logo">
      </a>
      <button class="hamburger" id="hamburger" aria-label="Toggle menu" aria-expanded="false">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <nav class="desktop-nav">
        <?php
        wp_nav_menu(array(
          'theme_location' => 'primary',
          'container' => false,
          'fallback_cb' => false,
        ));
        ?>
      </nav>
    </div>
  </header>
  <nav id="mobile-menu" class="mobile-menu">
    <?php
    wp_nav_menu(array(
      'theme_location' => 'primary',
      'container' => false,
      'fallback_cb' => false,
    ));
    ?>
  </nav>
  <div id="menu-overlay" class="menu-overlay"></div>
  <main>