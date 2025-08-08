// Google Analytics

// Initialize the dataLayer
window.dataLayer = window.dataLayer || [];

// Function to push events to the dataLayer
function gtag(){dataLayer.push(arguments);}

// Load the Google Tag Manager script asynchronously
(function() {
  var gtagScript = document.createElement('script');
  gtagScript.async = true;
  gtagScript.src = 'https://www.googletagmanager.com/gtag/js?id=UA-343131731';
  var firstScript = document.getElementsByTagName('script')[0];
  firstScript.parentNode.insertBefore(gtagScript, firstScript);
})();

// Configure Google Analytics with your Property ID
gtag('js', new Date());
gtag('config', 'UA-343131731');
