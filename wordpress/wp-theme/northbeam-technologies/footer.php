</main>
<footer class="site-footer">
  <div class="container footer-inner">
    <div class="footer-brand">
      <a class="footer-logo" href="<?php echo esc_url(home_url('/')); ?>">
        <img src="<?php echo esc_url(get_template_directory_uri() . '/assets/logo-dark.svg'); ?>" alt="<?php echo esc_attr(get_bloginfo('name')); ?>">
      </a>
      <p class="footer-tagline"><?php esc_html_e('Software and marketing built for global growth.', 'northbeam-technologies'); ?></p>
      <?php
      $northbeam_facebook  = trim((string) apply_filters('northbeam_facebook_url', get_theme_mod('northbeam_facebook_url', 'https://www.facebook.com/northbeam-technologies')));
      $northbeam_linkedin  = trim((string) apply_filters('northbeam_linkedin_url', get_theme_mod('northbeam_linkedin_url', 'https://www.linkedin.com/company/northbeam-technologies')));
      $northbeam_x         = trim((string) apply_filters('northbeam_x_url', get_theme_mod('northbeam_x_url', 'https://x.com/northbeam-technologies')));
      if ($northbeam_facebook !== '' || $northbeam_linkedin !== '' || $northbeam_x !== '') :
      ?>
      <div class="footer-social" aria-label="<?php esc_attr_e('Social media', 'northbeam-technologies'); ?>">
        <p class="footer-social-label"><?php esc_html_e('Follow us', 'northbeam-technologies'); ?></p>
        <div class="footer-social-links">
          <?php if ($northbeam_facebook !== '') : ?>
            <a class="footer-social-link" href="<?php echo esc_url($northbeam_facebook); ?>" target="_blank" rel="noopener noreferrer" aria-label="<?php esc_attr_e('Facebook', 'northbeam-technologies'); ?>">
              <?php echo wp_kses(northbeam_social_icon_facebook_svg(), northbeam_social_svg_kses_allowed()); ?>
            </a>
          <?php endif; ?>
          <?php if ($northbeam_linkedin !== '') : ?>
            <a class="footer-social-link" href="<?php echo esc_url($northbeam_linkedin); ?>" target="_blank" rel="noopener noreferrer" aria-label="<?php esc_attr_e('LinkedIn', 'northbeam-technologies'); ?>">
              <?php echo wp_kses(northbeam_social_icon_linkedin_svg(), northbeam_social_svg_kses_allowed()); ?>
            </a>
          <?php endif; ?>
          <?php if ($northbeam_x !== '') : ?>
            <a class="footer-social-link" href="<?php echo esc_url($northbeam_x); ?>" target="_blank" rel="noopener noreferrer" aria-label="<?php esc_attr_e('X', 'northbeam-technologies'); ?>">
              <?php echo wp_kses(northbeam_social_icon_x_svg(), northbeam_social_svg_kses_allowed()); ?>
            </a>
          <?php endif; ?>
        </div>
      </div>
      <?php endif; ?>
    </div>
    <div class="footer-nav-wrap">
      <nav class="footer-nav" aria-label="<?php esc_attr_e('Footer', 'northbeam-technologies'); ?>">
        <?php if (has_nav_menu('footer')) : ?>
          <?php
          wp_nav_menu(array(
            'theme_location' => 'footer',
            'container' => false,
            'menu_class' => 'footer-menu',
            'fallback_cb' => false,
            'depth' => 1,
          ));
          ?>
        <?php else : ?>
          <?php
          $northbeam_footer_links = apply_filters(
            'northbeam_footer_fallback_links',
            array(
              array(
                'url'   => home_url('/about-us/'),
                'label' => __('About Us', 'northbeam-technologies'),
              ),
              array(
                'url'   => home_url('/service/'),
                'label' => __('Our Services', 'northbeam-technologies'),
              ),
              array(
                'url'   => home_url('/project/'),
                'label' => __('Projects', 'northbeam-technologies'),
              ),
              array(
                'url'   => home_url('/insights/'),
                'label' => __('Insights', 'northbeam-technologies'),
              ),
              array(
                'url'   => home_url('/contact/'),
                'label' => __('Contact Us', 'northbeam-technologies'),
              ),
            )
          );
          ?>
          <ul class="footer-menu">
            <?php foreach ($northbeam_footer_links as $item) : ?>
              <li><a href="<?php echo esc_url($item['url']); ?>"><?php echo esc_html($item['label']); ?></a></li>
            <?php endforeach; ?>
          </ul>
        <?php endif; ?>
      </nav>
    </div>
  </div>
  <div class="footer-bottom">
    <div class="container footer-meta">
      <div class="footer-meta-row">
        <p class="footer-copyright">© <?php echo esc_html(date('Y')); ?> <?php echo esc_html(get_bloginfo('name')); ?>. <?php esc_html_e('All rights reserved.', 'northbeam-technologies'); ?></p>
        <?php
        $northbeam_footer_secondary = apply_filters(
          'northbeam_footer_secondary_links',
          array(
            array(
              'url'   => home_url('/privacy-policy/'),
              'label' => __('Privacy', 'northbeam-technologies'),
            ),
            array(
              'url'   => home_url('/terms/'),
              'label' => __('Terms', 'northbeam-technologies'),
            ),
            array(
              'url'   => home_url('/careers/'),
              'label' => __('Careers', 'northbeam-technologies'),
            ),
            array(
              'url'   => home_url('/sitemap/'),
              'label' => __('Sitemap', 'northbeam-technologies'),
            ),
            array(
              'url'   => home_url('/support/'),
              'label' => __('Support', 'northbeam-technologies'),
            ),
          )
        );
        if (!empty($northbeam_footer_secondary)) :
        ?>
        <nav class="footer-secondary" aria-label="<?php esc_attr_e('Legal and support', 'northbeam-technologies'); ?>">
          <ul class="footer-menu-secondary">
            <?php foreach ($northbeam_footer_secondary as $item) : ?>
              <li><a href="<?php echo esc_url($item['url']); ?>"><?php echo esc_html($item['label']); ?></a></li>
            <?php endforeach; ?>
          </ul>
        </nav>
        <?php endif; ?>
      </div>
    </div>
  </div>
</footer>
<?php wp_footer(); ?>
</body>
</html>
