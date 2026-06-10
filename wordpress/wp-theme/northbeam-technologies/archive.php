<?php get_header(); ?>
<div class="northbeam-inner northbeam-archive">
  <section class="archive-query-section">
    <div class="container">
      <header class="archive-header">
        <?php if (is_category()) : ?>
          <h1><?php single_cat_title(); ?></h1>
        <?php else : ?>
          <h1><?php the_archive_title(); ?></h1>
        <?php endif; ?>

        <?php if (get_the_archive_description()) : ?>
          <div class="archive-description"><?php the_archive_description(); ?></div>
        <?php endif; ?>
      </header>

      <?php if (have_posts()) : ?>
        <div class="archive-post-grid">
          <?php while (have_posts()) : the_post(); ?>
            <article <?php post_class('archive-post-card'); ?>>
              <?php if (has_post_thumbnail()) : ?>
                <a class="archive-post-card__image" href="<?php the_permalink(); ?>" aria-label="<?php the_title_attribute(); ?>">
                  <?php the_post_thumbnail('large'); ?>
                </a>
              <?php endif; ?>

              <div class="archive-post-card__content">
                <h2 class="archive-post-card__title">
                  <a href="<?php the_permalink(); ?>"><?php the_title(); ?></a>
                </h2>

                <div class="archive-post-card__excerpt">
                  <?php the_excerpt(); ?>
                </div>

                <div class="archive-post-card__meta">
                  <?php echo get_avatar(get_the_author_meta('ID'), 48); ?>
                  <a class="archive-post-card__author" href="<?php echo esc_url(get_author_posts_url(get_the_author_meta('ID'))); ?>">
                    <?php echo esc_html(get_the_author()); ?>
                  </a>
                  <span class="archive-post-card__divider" aria-hidden="true"></span>
                  <time datetime="<?php echo esc_attr(get_the_date(DATE_W3C)); ?>">
                    <?php echo esc_html(human_time_diff(get_the_time('U'), current_time('timestamp')) . ' ago'); ?>
                  </time>
                </div>
              </div>
            </article>
          <?php endwhile; ?>
        </div>

        <?php the_posts_pagination(array(
          'mid_size'  => 1,
          'prev_text' => __('Newer', 'northbeam-technologies'),
          'next_text' => __('Older', 'northbeam-technologies'),
        )); ?>
      <?php else : ?>
        <p class="archive-empty"><?php esc_html_e('No content found.', 'northbeam-technologies'); ?></p>
      <?php endif; ?>
    </div>
  </section>
</div>
<?php get_footer(); ?>
