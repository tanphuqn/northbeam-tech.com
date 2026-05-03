<?php get_header(); ?>
<section>
  <div class="container">
    <?php if (have_posts()) : while (have_posts()) : the_post(); ?>
      <article class="content-card">
        <h1><?php the_title(); ?></h1>
        <p><small><?php echo esc_html(get_the_date()); ?></small></p>
        <?php the_content(); ?>
      </article>
    <?php endwhile; endif; ?>
  </div>
</section>
<?php get_footer(); ?>
