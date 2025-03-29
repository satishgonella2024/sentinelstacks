// Wait for DOM to be fully loaded
document.addEventListener('DOMContentLoaded', function() {
    // Preloader
    setTimeout(function() {
        const preloader = document.querySelector('.preloader');
        if (preloader) {
            preloader.classList.add('loaded');
        }
    }, 500);

    // Mobile menu toggle
    const menuToggle = document.querySelector('.menu-toggle');
    const nav = document.querySelector('.nav');
    
    if (menuToggle && nav) {
        menuToggle.addEventListener('click', function() {
            nav.classList.toggle('active');
        });
    }

    // Smooth scroll for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            const href = this.getAttribute('href');
            
            if (href !== '#') {
                e.preventDefault();
                
                const targetElement = document.querySelector(href);
                if (targetElement) {
                    // Close mobile menu if open
                    if (nav.classList.contains('active')) {
                        nav.classList.remove('active');
                    }
                    
                    // Scroll to target
                    window.scrollTo({
                        top: targetElement.offsetTop - 80, // Account for header height
                        behavior: 'smooth'
                    });
                }
            }
        });
    });

    // Animate elements on scroll
    const animateElements = document.querySelectorAll('.feature-card, .step, .example-card, .download-option');
    
    // Check if IntersectionObserver is supported
    if ('IntersectionObserver' in window) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate-fadeIn');
                    observer.unobserve(entry.target);
                }
            });
        }, {
            threshold: 0.1
        });

        animateElements.forEach(element => {
            observer.observe(element);
        });
    } else {
        // Fallback for browsers that don't support IntersectionObserver
        animateElements.forEach(element => {
            element.classList.add('animate-fadeIn');
        });
    }

    // Newsletter form submission
    const newsletterForm = document.querySelector('.newsletter-form');
    
    if (newsletterForm) {
        newsletterForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const emailInput = this.querySelector('input[type="email"]');
            const email = emailInput.value.trim();
            
            if (email) {
                // In a real implementation, you would send this to your backend
                console.log('Newsletter subscription for:', email);
                
                // Show success message
                const formContainer = this.parentElement;
                const successMessage = document.createElement('div');
                successMessage.className = 'newsletter-success';
                successMessage.innerHTML = '<p>Thank you for subscribing! We\'ll keep you updated on SentinelStacks news.</p>';
                
                formContainer.innerHTML = '';
                formContainer.appendChild(successMessage);
            }
        });
    }

    // Code block copy functionality
    document.querySelectorAll('.code-block').forEach(codeBlock => {
        const code = codeBlock.querySelector('code');
        
        if (code) {
            const copyButton = document.createElement('button');
            copyButton.className = 'copy-button';
            copyButton.innerHTML = '<i class="fas fa-copy"></i>';
            copyButton.title = 'Copy to clipboard';
            
            codeBlock.style.position = 'relative';
            copyButton.style.position = 'absolute';
            copyButton.style.top = '8px';
            copyButton.style.right = '8px';
            copyButton.style.background = 'rgba(255, 255, 255, 0.1)';
            copyButton.style.border = 'none';
            copyButton.style.borderRadius = '4px';
            copyButton.style.padding = '4px 8px';
            copyButton.style.color = 'white';
            copyButton.style.cursor = 'pointer';
            
            codeBlock.appendChild(copyButton);
            
            copyButton.addEventListener('click', function() {
                const textToCopy = code.textContent;
                
                navigator.clipboard.writeText(textToCopy).then(() => {
                    copyButton.innerHTML = '<i class="fas fa-check"></i>';
                    
                    setTimeout(() => {
                        copyButton.innerHTML = '<i class="fas fa-copy"></i>';
                    }, 2000);
                }).catch(err => {
                    console.error('Could not copy text: ', err);
                });
            });
        }
    });
});
