/**
 * Generates a complete theme based on a few base colors
 * @param {Object} options - The base colors and options for the theme
 * @param {string} options.primary - Primary color (for buttons, active elements)
 * @param {string} options.secondary - Secondary color (for secondary buttons)
 * @param {string} options.background - Base background color
 * @param {string} [options.text] - Main text color (will be auto-generated if not provided)
 * @param {string} [options.accent] - Accent color (will default to a variation of primary if not provided)
 * @param {boolean} [options.darkMode] - Whether to generate a dark mode theme. If not provided, will be auto-detected based on background color brightness
 * @returns {Object} Complete theme in the same format as the JSON file
 */
function generateTheme(options) {
  const {
    primary,
    secondary,
    background,
    text: customText,
    accent: customAccent,
    darkMode: userSpecifiedDarkMode
  } = options;

  // Auto-detect if dark mode is appropriate based on background color brightness
  const isDarkBackground = (color) => {
    const rgb = hexToRgb(color);
    const brightness = (rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000;
    return brightness < 125; // Lower threshold indicates a dark color
  };

  // Use user-specified dark mode setting if provided, otherwise auto-detect
  const darkMode = userSpecifiedDarkMode !== undefined 
    ? userSpecifiedDarkMode 
    : isDarkBackground(background);

  // Helper functions for color manipulation
  const hexToRgb = (hex) => {
    const shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
    const fullHex = hex.replace(shorthandRegex, (m, r, g, b) => r + r + g + g + b + b);
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(fullHex);
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16)
        }
      : null;
  };

  const rgbToHex = (r, g, b) => {
    return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
  };

  const rgbToString = (rgb, alpha = 1) => {
    return alpha < 1
      ? `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${alpha})`
      : `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})`;
  };

  // Blend two colors
  const blendColors = (color1, color2, ratio = 0.5) => {
    const rgb1 = hexToRgb(color1);
    const rgb2 = hexToRgb(color2);
    
    const r = Math.round(rgb1.r * (1 - ratio) + rgb2.r * ratio);
    const g = Math.round(rgb1.g * (1 - ratio) + rgb2.g * ratio);
    const b = Math.round(rgb1.b * (1 - ratio) + rgb2.b * ratio);
    
    return rgbToString({ r, g, b });
  };

  // Lighten or darken a color
  const adjustBrightness = (color, factor) => {
    const rgb = hexToRgb(color);
    const adjust = (value) => Math.min(255, Math.max(0, Math.round(value * factor)));
    
    return rgbToString({
      r: adjust(rgb.r),
      g: adjust(rgb.g),
      b: adjust(rgb.b)
    });
  };

  // Calculate contrast color (for text)
  const getContrastColor = (bgColor) => {
    const rgb = hexToRgb(bgColor);
    const brightness = (rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000;
    return brightness > 128 ? 'rgb(50, 50, 50)' : 'rgb(255, 255, 255)';
  };

  // Get a color that ensures readability against the background
  const getReadableTextColor = (bgColor, darkPreferred = false) => {
    const rgb = hexToRgb(bgColor);
    const brightness = (rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000;
    
    if (darkPreferred) {
      return brightness > 200 ? 'rgb(50, 50, 50)' : 'rgb(83, 83, 83)';
    } else {
      return brightness > 150 ? 'rgb(50, 50, 50)' : 'rgb(255, 255, 255)';
    }
  };

  // Convert colors to RGB objects
  const primaryRgb = hexToRgb(primary);
  const secondaryRgb = hexToRgb(secondary);
  const backgroundRgb = hexToRgb(background);
  
  // Determine text color based on background if not provided
  const text = customText || getReadableTextColor(background, !darkMode);
  
  // Generate accent color if not provided
  const accent = customAccent || blendColors(primary, secondary, 0.3);
  
  // Generate hover and active variations
  const primaryHover = adjustBrightness(primary, darkMode ? 0.85 : 1.15);
  const primaryActive = adjustBrightness(primary, darkMode ? 0.7 : 1.3);
  
  const secondaryHover = adjustBrightness(secondary, darkMode ? 0.9 : 1.05);
  const secondaryActive = adjustBrightness(secondary, darkMode ? 0.8 : 1.1);
  
  // Text colors for buttons
  const primaryTextColor = getContrastColor(primary);
  const secondaryTextColor = getContrastColor(secondary);
  
  // Calculate contrast ratio between two colors (simplified version)
  const calculateContrastRatio = (color1, color2) => {
    const rgb1 = typeof color1 === 'string' ? hexToRgb(color1) : color1;
    const rgb2 = typeof color2 === 'string' ? hexToRgb(color2) : color2;
    
    // Calculate relative luminance
    const getLuminance = (rgb) => {
      const rsrgb = rgb.r / 255;
      const gsrgb = rgb.g / 255;
      const bsrgb = rgb.b / 255;
      
      const r = rsrgb <= 0.03928 ? rsrgb / 12.92 : Math.pow((rsrgb + 0.055) / 1.055, 2.4);
      const g = gsrgb <= 0.03928 ? gsrgb / 12.92 : Math.pow((gsrgb + 0.055) / 1.055, 2.4);
      const b = bsrgb <= 0.03928 ? bsrgb / 12.92 : Math.pow((bsrgb + 0.055) / 1.055, 2.4);
      
      return 0.2126 * r + 0.7152 * g + 0.0722 * b;
    };
    
    const l1 = getLuminance(rgb1);
    const l2 = getLuminance(rgb2);
    
    // Return contrast ratio
    return (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
  };
  
  // Adjust link color to ensure readability against background
  const getReadableLinkColor = (primaryColor, backgroundColor) => {
    const minContrastRatio = 4.5; // WCAG AA minimum for normal text
    let linkColor = primaryColor;
    const primaryRgb = typeof primaryColor === 'string' ? hexToRgb(primaryColor) : primaryColor;
    
    // Check if primary color has enough contrast with background
    if (calculateContrastRatio(primaryRgb, backgroundRgb) < minContrastRatio) {
      // If we're in dark mode, brighten the link color
      if (darkMode) {
        // Start with primary and gradually brighten until we hit desired contrast
        let factor = 1.2;
        const brightenedRgb = { ...primaryRgb };
        
        while (calculateContrastRatio(brightenedRgb, backgroundRgb) < minContrastRatio && factor <= 2.0) {
          brightenedRgb.r = Math.min(255, Math.round(primaryRgb.r * factor));
          brightenedRgb.g = Math.min(255, Math.round(primaryRgb.g * factor));
          brightenedRgb.b = Math.min(255, Math.round(primaryRgb.b * factor));
          factor += 0.1;
        }
        
        linkColor = rgbToString(brightenedRgb);
      } else {
        // For light modes, darken the link color
        let factor = 0.8;
        const darkenedRgb = { ...primaryRgb };
        
        while (calculateContrastRatio(darkenedRgb, backgroundRgb) < minContrastRatio && factor >= 0.3) {
          darkenedRgb.r = Math.max(0, Math.round(primaryRgb.r * factor));
          darkenedRgb.g = Math.max(0, Math.round(primaryRgb.g * factor));
          darkenedRgb.b = Math.max(0, Math.round(primaryRgb.b * factor));
          factor -= 0.1;
        }
        
        linkColor = rgbToString(darkenedRgb);
      }
    }
    
    return linkColor;
  };
  
  // Get readable link colors
  const linkColor = getReadableLinkColor(primaryRgb, backgroundRgb);
  const linkRgb = typeof linkColor === 'string' ? hexToRgb(linkColor) : linkColor;
  const linkHover = adjustBrightness(rgbToString(linkRgb), darkMode ? 1.2 : 0.8);
  const linkActive = adjustBrightness(rgbToString(linkRgb), darkMode ? 1.4 : 0.6);

  // Generate the theme
  return {
    "links": {
      "default": typeof linkColor === 'string' ? linkColor : rgbToString(linkColor),
      "hover": linkHover,
      "active": linkActive,
      "disabled": "rgba(190, 190, 190, 0.6)"
    },
    "buttons": {
      "primary": {
        "default": {
          "background": rgbToString(primaryRgb),
          "text": primaryTextColor
        },
        "hover": {
          "background": primaryHover,
          "text": primaryTextColor
        },
        "active": {
          "background": primaryActive,
          "text": primaryTextColor
        },
        "disabled": {
          "background": "rgba(120, 120, 120, 0.5)",
          "text": "rgba(255, 255, 255, 0.6)"
        }
      },
      "secondary": {
        "default": {
          "background": rgbToString(secondaryRgb),
          "text": secondaryTextColor
        },
        "hover": {
          "background": secondaryHover,
          "text": secondaryTextColor
        },
        "active": {
          "background": secondaryActive,
          "text": secondaryTextColor
        },
        "disabled": {
          "background": "rgba(120, 120, 120, 0.5)",
          "text": "rgba(255, 255, 255, 0.6)"
        }
      }
    },
    "uploader": {
      "background": rgbToString(backgroundRgb, 0.769),
      "text": text,
      "header": {
        "background": rgbToString(secondaryRgb, 0.5),
        "text": text
      },
      "item": {
        "background": adjustBrightness(background, 1.05),
        "text": text
      }
    },
    "panel": {
      "background": rgbToString(backgroundRgb),
      "text": text,
      "item": {
        "background": adjustBrightness(background, darkMode ? 0.9 : 0.95),
        "text": text,
        "hover": {
          "background": adjustBrightness(background, darkMode ? 0.8 : 0.9),
          "text": text
        },
        "action": {
          "hover": {
            "backgroundFrom": rgbToString(secondaryRgb, 0.42),
            "backgroundTo": secondaryHover
          }
        }
      },
      "header": {
        "background": rgbToString(primaryRgb, 0.8),
        "text": primaryTextColor
      },
      "textAlt": adjustBrightness(text, 1.4),
      "subheader": {
        "background": `linear-gradient(to bottom, ${rgbToString(backgroundRgb)} 0%, ${adjustBrightness(background, darkMode ? 0.9 : 0.95)} 100%)`,
        "text": text
      },
      "nav": {
        "item": {
          "background": adjustBrightness(background, darkMode ? 0.9 : 0.95),
          "text": text,
          "hover": {
            "background": rgbToString(primaryRgb, 0.8),
            "text": primaryTextColor
          }
        }
      },
      "section": {
        "background": adjustBrightness(background, darkMode ? 0.9 : 0.95),
        "text": text,
        "alt": {
          "background": adjustBrightness(background, darkMode ? 0.8 : 0.9),
          "text": text
        }
      }
    },
    "tabs": {
      "bar": {
        "background": adjustBrightness(background, darkMode ? 0.8 : 0.9)
      },
      "tab": {
        "background": adjustBrightness(background, darkMode ? 0.9 : 0.95),
        "text": text,
        "hover": {
          "background": rgbToString(primaryRgb, 0.8),
          "text": primaryTextColor
        },
        "active": {
          "background": rgbToString(primaryRgb, 0.8),
          "text": primaryTextColor
        }
      }
    },
    "table": {
      "background": adjustBrightness(background, darkMode ? 0.8 : 0.9),
      "text": text,
      "header": {
        "background": rgbToString(primaryRgb),
        "text": primaryTextColor
      },
      "row": {
        "background": rgbToString(backgroundRgb),
        "text": text,
        "alt": {
          "background": adjustBrightness(background, darkMode ? 0.9 : 0.95),
          "text": text
        }
      }
    },
    "input": {
      "background": adjustBrightness(background, darkMode ? 0.9 : 1.05),
      "text": text,
      "border": {
        "default": primaryActive,
        "hover": primaryHover,
        "focus": rgbToString(primaryRgb)
      },
      "placeholder": adjustBrightness(text, 1.7)
    },
    "checkbox": {
      "background": adjustBrightness(background, darkMode ? 0.9 : 1.05),
      "checked": {
        "background": rgbToString(primaryRgb),
        "check": primaryTextColor
      }
    },
    "label": {
      "text": text
    },
    "progressBar": {
      "background": adjustBrightness(background, darkMode ? 0.8 : 0.9),
      "fill": rgbToString(primaryRgb),
      "text": text
    },
    "overlay": {
      "background": rgbToString(backgroundRgb, 0.2)
    },
    "dimensions": {
      "button": {
        "height": "50px",
        "width": "100%",
        "iconOnlyWidth": "50px",
        "borderRadius": "10px"
      },
      "panel": {
        "borderRadius": "10px"
      },
      "tabs": {
        "borderRadius": "5px 5px 0 0"
      },
      "settings": {
        "width": "100vw",
        "height": "100vh"
      }
    }
  };
}

// Example usage:
// const theme = generateTheme({
//   primary: '#DB7FA3',    // Pink
//   secondary: '#FCF5F8',  // Light pink
//   background: '#FCF0F5', // Very light pink
//   darkMode: false
// });

// Utility to convert CSS variables to theme JSON
function cssVarsToThemeJson(cssVarsText) {
  // Parse CSS variables into an object
  const cssVars = {};
  cssVarsText.split('\n').forEach(line => {
    const match = line.match(/--([^:]+):\s*(.*);?$/);
    if (match) {
      cssVars[match[1].trim()] = match[2].trim();
    }
  });
  
  // Create the JSON structure
  return {
    "links": {
      "default": cssVars["link-color"],
      "hover": cssVars["link-color-hover"],
      "active": cssVars["link-color-active"],
      "disabled": cssVars["link-color-disabled"]
    },
    "buttons": {
      "primary": {
        "default": {
          "background": cssVars["primary-button-background-color"],
          "text": cssVars["primary-button-text-color"]
        },
        "hover": {
          "background": cssVars["primary-button-background-color-hover"],
          "text": cssVars["primary-button-text-color-hover"]
        },
        "active": {
          "background": cssVars["primary-button-background-color-hover"],
          "text": cssVars["primary-button-text-color-hover"]
        },
        "disabled": {
          "background": cssVars["primary-button-background-color-disabled"],
          "text": cssVars["primary-button-text-color-disabled"]
        }
      },
      "secondary": {
        "default": {
          "background": cssVars["secondary-button-background-color"],
          "text": cssVars["secondary-button-text-color"]
        },
        "hover": {
          "background": cssVars["secondary-button-background-color-hover"],
          "text": cssVars["secondary-button-text-color-hover"]
        },
        "active": {
          "background": cssVars["secondary-button-background-color-hover"],
          "text": cssVars["secondary-button-text-color-hover"]
        },
        "disabled": {
          "background": cssVars["secondary-button-background-color-disabled"],
          "text": cssVars["secondary-button-text-color-disabled"]
        }
      }
    },
    // ... Continue with all other properties
  };
}

// Utility to convert theme JSON to CSS variables format
function themeJsonToCssVars(themeJson) {
  let cssVars = '';
  
  // Links
  cssVars += `/* link colors */\n`;
  cssVars += `--link-color: ${themeJson.links.default};\n`;
  cssVars += `--link-color-hover: ${themeJson.links.hover};\n`;
  cssVars += `--link-color-active: ${themeJson.links.active};\n`;
  cssVars += `--link-color-disabled: ${themeJson.links.disabled};\n\n`;
  
  // Buttons
  cssVars += `/* Button Colors */\n`;
  cssVars += `--primary-button-background-color: ${themeJson.buttons.primary.default.background};\n`;
  cssVars += `--primary-button-text-color: ${themeJson.buttons.primary.default.text};\n\n`;
  cssVars += `--primary-button-background-color-hover: ${themeJson.buttons.primary.hover.background};\n`;
  cssVars += `--primary-button-text-color-hover: ${themeJson.buttons.primary.hover.text};\n\n`;
  cssVars += `--primary-button-background-color-disabled: ${themeJson.buttons.primary.disabled.background};\n`;
  cssVars += `--primary-button-text-color-disabled: ${themeJson.buttons.primary.disabled.text};\n\n`;
  cssVars += `--secondary-button-background-color: ${themeJson.buttons.secondary.default.background};\n`;
  cssVars += `--secondary-button-text-color: ${themeJson.buttons.secondary.default.text};\n\n`;
  cssVars += `--secondary-button-background-color-hover: ${themeJson.buttons.secondary.hover.background};\n`;
  cssVars += `--secondary-button-text-color-hover: ${themeJson.buttons.secondary.hover.text};\n\n`;
  cssVars += `--secondary-button-background-color-disabled: ${themeJson.buttons.secondary.disabled.background};\n`;
  cssVars += `--secondary-button-text-color-disabled: ${themeJson.buttons.secondary.disabled.text};\n\n`;
  
  // ... Continue with all other properties
  
  return cssVars;
}

export { generateTheme, cssVarsToThemeJson, themeJsonToCssVars }