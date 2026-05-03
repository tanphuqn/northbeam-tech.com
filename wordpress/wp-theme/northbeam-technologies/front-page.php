<?php
/**
 * Front page: render block content from the assigned static homepage (paste from wp-pages/home.html).
 * Matches static-html/index.html when the same blocks are used.
 */
get_header();
?>
<div class="northbeam-front">
  <?php
  if (have_posts()) :
    while (have_posts()) :
      the_post();
      the_content();
    endwhile;
  endif;
  ?>
</div>
<?php get_footer(); ?>
