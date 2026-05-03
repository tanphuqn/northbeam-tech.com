<?php
/**
 * Inner pages: full-width block content (no extra .container/.content-card).
 * Same block styles as home (.northbeam-inner mirrors .northbeam-front in style.css).
 * Page title is not shown here — use the H1 inside your pasted blocks (or add one).
 */
get_header();
?>
<?php if (have_posts()) : while (have_posts()) : the_post(); ?>
  <?php if (is_front_page()) : ?>
    <div class="northbeam-front">
      <?php the_content(); ?>
    </div>
  <?php else : ?>
    <div class="northbeam-inner">
      <?php the_content(); ?>
    </div>
  <?php endif; ?>
<?php endwhile; endif; ?>
<?php get_footer(); ?>
