<?php get_header(); ?>
<section>
  <div class="container">
    <h1>Latest Posts</h1>
    <?php if (have_posts()) : while (have_posts()) : the_post(); ?>
      <article class="content-card">
        <h2><a href="<?php the_permalink(); ?>"><?php the_title(); ?></a></h2>
        <p><small><?php echo esc_html(get_the_date()); ?></small></p>
        <?php the_excerpt(); ?>
      </article>
    <?php endwhile; else : ?>
      <article class="content-card">
        <p>No posts found.</p>
      </article>
    <?php endif; ?>
  </div>
</section>
<?php get_footer(); ?>
