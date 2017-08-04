(function () {
  var tracker = document.getElementById('weatherglass-tracker');

  var src = tracker.getAttribute('src');
  src = src.replace('/track.js', '/track.gif');

  var site_id = tracker.getAttribute('data-site-id');
  var resource = document.location.href;
  var referrer = document.referrer;
  var title = document.title;
  var user_agent = navigator.userAgent;
 
  var url = new String(src);
  url += "?site_id=" + site_id;
  url += "&resource=" + resource;
  url += "&referrer=" + referrer;
  url += "&title=" + title;
  url += "&user_agent=" + user_agent;

  var i = new Image();
  i.src = url;
})();
