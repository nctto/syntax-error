document.addEventListener('DOMContentLoaded', function () {
  window.addEventListener('scroll', function (e) {
    const scrollY = window.scrollY;
    const nav = document.getElementById('nav');
    const navTop = nav.offsetTop;
    const toolbar = document.getElementById('toolbar');
    console.log(toolbar.offsetTop);
    const toolbarTop = toolbar.offsetTop - 50;
    const projects = document.getElementById('projects');
    const projectsTop = projects.offsetTop;
    const projectsWidth = projects.offsetWidth - 10;
    const tclass = ['fixed', 'top-3', 'z-10'];
    if (scrollY < projectsTop) {
      tclass.forEach((c) => {
        toolbar.classList.remove(c);
      });
      toolbar.style.width = '';
      // nav.classList.remove('bg-white');
    } else {
      tclass.forEach((c) => {
        toolbar.classList.add(c);
      });
      toolbar.style.width = `${projectsWidth}px`;
      // nav.classList.add('bg-white');
    }
  });
});
