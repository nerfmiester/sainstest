package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var mockDoc = `
<!DOCTYPE html>
<html class="noJs" xmlns:wairole="http://www.w3.org/2005/01/wai-rdf/GUIRoleTaxonomy#" xmlns:waistate="http://www.w3.org/2005/07/aaa" lang="en" xml:lang="en">
<!-- BEGIN CategoriesDisplay.jsp -->
<head>
    <title>Ripe &amp; ready | Sainsbury&#039;s</title>
    <meta name="description" content="Buy Ripe &amp; ready online from Sainsbury&#039;s, the same great quality, freshness and choice you&#039;d find in store. Choose from 1 hour delivery slots and collect Nectar points."/>
    <meta name="keyword" content=""/>
    <link rel="canonical" href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/ripe---ready" />

        <meta content="NOINDEX, FOLLOW" name="ROBOTS" />
    <!-- BEGIN CommonCSSToInclude.jspf --><!--[if IE 8]>
    <link type="text/css" href="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/css/main-ie8.min.css" rel="stylesheet" media="all" />
	<![endif]-->

    <!--[if !IE 8]><!-->
    <link type="text/css" href="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/css/main.min.css" rel="stylesheet" media="all" />
    <!--<![endif]-->


	<link type="text/css" href="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/wcassets/groceries/css/espot.css" rel="stylesheet" media="all" />
	<link type="text/css" href="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/css/print.css" rel="stylesheet" media="print"/>
<!-- END CommonCSSToInclude.jspf --><!-- BEGIN CommonJSToInclude.jspf -->
<meta name="CommerceSearch" content="storeId_10151" />



<script type="text/javascript">
    var _deliverySlotInfo = {
            expiryDateTime: '',
            currentDateTime: 'November 25,2015 16:53:57',
            ajaxCountDownUrl: 'CountdownDisplayView?langId=44&storeId=10151',
            ajaxExpiredUrl: 'DeliverySlotExpiredDisplayView?langId=44&storeId=10151&currentPageUrl=http%3a%2f%2fwww.sainsburys.co.uk%2fwebapp%2fwcs%2fstores%2fservlet%2fCategoryDisplay%3fmsg%3d%26categoryId%3d185749%26langId%3d44%26storeId%3d10151%26krypto%3ddwlvaeB6%252FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%250A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%252F%252BHeNnUqybiZXu%252FL47P9A658zhrWd08mA5Y%250Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%252BardwWtMA49XQA4Iqwf%252BSvFr8RJOHK%250Afp2%252Fk0F6LH6%252Fmq5%252FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%252FydqEDvTdI5qgO6sKl0Q%253D&AJAXCall=true'
        }
    var _amendOrderSlotInfo = {
            expiryDateTime: '',
            currentDateTime: 'November 25,2015 16:53:57',
            ajaxAmendOrderExpiryUrl: 'AjaxOrderAmendSlotExpiryView?langId=44&storeId=10151&currentPageUrl=http%3a%2f%2fwww.sainsburys.co.uk%2fwebapp%2fwcs%2fstores%2fservlet%2fCategoryDisplay%3fmsg%3d%26categoryId%3d185749%26langId%3d44%26storeId%3d10151%26krypto%3ddwlvaeB6%252FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%250A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%252F%252BHeNnUqybiZXu%252FL47P9A658zhrWd08mA5Y%250Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%252BardwWtMA49XQA4Iqwf%252BSvFr8RJOHK%250Afp2%252Fk0F6LH6%252Fmq5%252FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%252FydqEDvTdI5qgO6sKl0Q%253D'
        }
    var _commonPageInfo = {
        currentUrl: 'http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D',
        storeId: '10151',
        langId: '44'
    }
</script>

        <script type="text/javascript">
	    var _rhsCheckPostCodeRuleset = {
	          postCode: {
	                isEmpty: {
	                      param: true,
	                      text: 'Sorry, this postcode has not been recognised - Please try again.',
	                      msgPlacement: "#checkPostCodePanel #Rhs_checkPostCode .field",
	                      elemToAddErrorClassTo: "#checkPostCodePanel #Rhs_checkPostCode .field"
	                },
	                minLength: {
	                      param: 5,
	                      text: 'Sorry, this entry must be at least 5 characters long.',
	                      msgPlacement: "#checkPostCodePanel #Rhs_checkPostCode .field",
	                      elemToAddErrorClassTo: "#checkPostCodePanel #Rhs_checkPostCode .field"
	                },
	                maxLength: {
	                      param: 8,
	                      text: 'Sorry, this postcode has not been recognised - Please try again.',
	                      msgPlacement: "#checkPostCodePanel #Rhs_checkPostCode .field",
	                      elemToAddErrorClassTo: "#checkPostCodePanel #Rhs_checkPostCode .field"
	                },
	                isPostcode: {
	                      param: true,
	                      text: 'Sorry, this postcode has not been recognised - Please try again.',
	                      msgPlacement: "#checkPostCodePanel #Rhs_checkPostCode .field",
	                      elemToAddErrorClassTo: "#checkPostCodePanel #Rhs_checkPostCode .field"
	                }
	          }
	    }
	    </script>

        <script type="text/javascript">
	    var _rhsLoginValidationRuleset = {
	        logonId: {
	            isEmpty: {
	                param: true,
	                text: 'Please enter your username in the space provided.',
	                msgPlacement: "fieldUsername",
	                elemToAddErrorClassTo: "fieldUsername"
	            },
	            notMatches: {
	                param: "#logonPassword",
	                text: 'Sorry, your details have not been recognised. Please try again.',
	                msgPlacement: "fieldUsername",
	                elemToAddErrorClassTo: "fieldUsername"
	            }
	        },
	        logonPassword: {
	            isEmpty: {
	                param: true,
	                text: 'Please enter your password in the space provided.',
	                msgPlacement: "fieldPassword",
	                elemToAddErrorClassTo: "fieldPassword"
	            },
	            minLength: {
	                param: 6,
	                text: 'Please enter your password in the space provided.',
	                msgPlacement: "fieldPassword",
	                elemToAddErrorClassTo: "fieldPassword"
	            }
	        }
	    }
	    </script>

<script type="text/javascript">
      var typeAheadTrigger = 2;
</script>

<!--<script type="text/javascript" data-dojo-config="isDebug: false, useCommentedJson: true,locale: 'en-gb', parseOnLoad: true, dojoBlankHtmlUrl:'/wcsstore/SainsburysStorefrontAssetStore/js/dojo.1.7.1/blank.html'" src="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/js/dojo.1.7.1/dojo/dojo.js"></script>-->




<script type="text/javascript" src="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/js/sainsburys.js"></script>


<script type="text/javascript">require(["dojo/parser", "dijit/layout/AccordionContainer", "dijit/layout/ContentPane", "dojox/widget/AutoRotator", "dojox/widget/rotator/Fade"]);</script>
<script type="text/javascript" src="http://c1.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/wcassets/groceries/scripts/page/faq.js"></script>


    <style id="antiCJ">.js body{display:none !important;}</style>
    <script type="text/javascript">if (self === top) {var antiCJ = document.getElementById("antiCJ");antiCJ.parentNode.removeChild(antiCJ);} else {top.location = self.location;}</script>
<!-- END CommonJSToInclude.jspf -->
</head>

<body id="shelfPage" class="shelfPage">
<div id="page">
    <!-- BEGIN StoreCommonUtilities.jspf --><!-- END StoreCommonUtilities.jspf --><!-- Header Nav Start --><!-- BEGIN LayoutContainerTop.jspf --><!-- BEGIN HeaderDisplay.jspf --><!-- BEGIN CachedHeaderDisplay.jsp -->

<ul id="skipLinks">
    <li><a href="#content">Skip to main content</a></li>
    <li><a href="#groceriesNav">Skip to groceries navigation menu</a></li>

</ul>

<div id="globalHeaderContainer">
    <div class="header globalHeader" id="globalHeader">
        <div class="globalNav">
	<ul>
		<li>
			<a href="http://www.sainsburys.co.uk">
			    <span class="moreSainsburysIcon"></span>
                Explore more at Sainsburys.co.uk
			</a>
		</li>
		<li>
			<a href="http://help.sainsburys.co.uk" rel="external">
			    <span class="helpCenterIcon"></span>
                Help Centre
			</a>
		</li>
		<li>
			<a href="http://stores.sainsburys.co.uk">
			    <span class="storeLocatorIcon"></span>
                Store Locator
			</a>
		</li>
		<li class="loginRegister">

					<a href="https://www.sainsburys.co.uk/sol/my_account/accounts_home.jsp">
						<span class="userIcon"></span>
                        Log in / Register
					</a>

		</li>
	</ul>
</div>

	    <div class="globalHeaderLogoSearch">
	        <!-- BEGIN LogoSearchNavBar.jspf -->

<a href="http://www.sainsburys.co.uk/shop/gb/groceries" class="mainLogo"><img src="http://www.sainsburys.co.uk/wcsstore/SainsburysStorefrontAssetStore/img/logo.png" alt="Sainsbury's" /></a>
<div class="searchBox" role="search">


    <form name="sol_search" method="get" action="SearchDisplay" id="globalSearchForm">

        <input type="hidden" name="viewTaskName" value="CategoryDisplayView" />
        <input type="hidden" name="recipesSearch" value="true" />
        <input type="hidden" name="orderBy" value="RELEVANCE" />


              <input type="hidden" name="skipToTrollyDisplay" value="false"/>

              <input type="hidden" name="favouritesSelection" value="0"/>

              <input type="hidden" name="level" value="2"/>

              <input type="hidden" name="langId" value="44"/>

              <input type="hidden" name="storeId" value="10151"/>


        <label for="search" class="access">Search for products</label>
        <input type="search" name="searchTerm" id="search" maxlength="150" value="" autocomplete="off" placeholder="Search" />
        <button type="button" id="clearSearch" class="clearSearch hidden">Clear the search field</button>
        <input type="submit" name="searchSubmit" id="searchSubmit" value="Search" />
    </form>

    <a class="findProduct" href="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/ShoppingListDisplay?catalogId=10122&action=ShoppingListDisplay&urlLangId=&langId=44&storeId=10151">Search for multiple products</a>
    <!-- ul class="searchNav">
        <li class="shoppingListLink"><a href="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/ShoppingListDisplay?catalogId=10122&action=ShoppingListDisplay&urlLangId=&langId=44&storeId=10151">Find a list of products</a></li>
        <li><a href="http://stores.sainsburys.co.uk">Store Locator</a></li>
        <li><a href="https://www.sainsburys.co.uk/sol/my_account/accounts_home.jsp">My Account</a></li>

                 <li><a href="https://www.sainsburys.co.uk/webapp/wcs/stores/servlet/QuickRegistrationFormView?catalogId=10122&amp;langId=44&amp;storeId=10151" >Register</a></li>

    </ul-->

</div>
<!-- END LogoSearchNavBar.jspf -->
        </div>
        <div id="groceriesNav" class="groceriesNav">
            <ul class="mainNav">
                <li>

                            <a class="active" href="http://www.sainsburys.co.uk/shop/gb/groceries"><strong>Groceries</strong></a>

                </li>
                <li>

                           <a href="http://www.sainsburys.co.uk/shop/gb/groceries/favourites">Favourites</a>

                </li>
                <li>

                          <a href="http://www.sainsburys.co.uk/shop/gb/groceries/great-offers">Great Offers</a>

                </li>
                <li>

                           <a href="http://www.sainsburys.co.uk/shop/gb/groceries/ideas-recipes">Ideas &amp; Recipes</a>

                </li>
                <li>

                           <a href="http://www.sainsburys.co.uk/shop/gb/groceries/benefits">Benefits</a>

                </li>
            </ul>
            <hr />

                    <p class="access">Groceries Categories</p>

                    <div class="subNav">
                        <ul>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/Christmas">Christmas</a>

                                   </li>

                                <li>

                                            <a class="active" href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg"><strong>Fruit &amp; veg</strong></a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/meat-fish">Meat &amp; fish</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/dairy-eggs-chilled">Dairy, eggs &amp; chilled</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/bakery">Bakery</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/frozen-">Frozen</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/food-cupboard">Food cupboard</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/drinks">Drinks</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/health-beauty">Health &amp; beauty</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/baby">Baby</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/household">Household</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/pet">Pet</a>

                                   </li>

                                <li>

                                            <a href="http://www.sainsburys.co.uk/shop/gb/groceries/home-ents">Home</a>

                                   </li>

                        </ul>
                    </div>

        </div>
    </div>
</div>
<!-- Generated on: Wed Nov 25 16:53:57 GMT 2015  -->
<!-- END CachedHeaderDisplay.jsp --><!-- END HeaderDisplay.jspf --><!-- END LayoutContainerTop.jspf --><!-- Header Nav End --><!-- Main Area Start -->
    <div id="main">
        <!-- Content Start -->
        <div class="article" id="content">

                  <div class="nav breadcrumb" id="breadcrumbNav">
                    <p class="access">You are here:</p>
                    <ul>

<li class="first"><span class="corner"></span><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg"><span>Fruit & veg</span></a>

        <span class="arrow"></span>

    <div>
        <p>Select an option:</p>
        <ul>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/great-prices-on-fruit---veg">Great prices on fruit & veg</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/flowers---seeds">Flowers & plants</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/new-in-season">In season</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-fruit">Fresh fruit</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-vegetables">Fresh vegetables</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-salad">Fresh salad</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-herbs-ingredients">Fresh herbs & ingredients</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/prepared-ready-to-eat">Prepared fruit, veg & salad</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/organic">Organic</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/taste-the-difference-185761-44">Taste the Difference</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fruit-veg-fairtrade">Fairtrade</a></li>

                <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/christmas-fruit---nut">Christmas fruit & nut</a></li>

        </ul>
    </div>
</li>

            <li class="second"><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-fruit"><span>Fresh fruit</span></a> <span class="arrow"></span>
                <div>
                <p>Select an option:</p>
                    <ul>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/all-fruit">All fruit</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/ripe---ready">Ripe & ready</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/bananas-grapes">Bananas & grapes</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/apples-pears-rhubarb">Apples, pears & rhubarb</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/berries-cherries-currants">Berries, cherries & currants</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/citrus-fruit">Citrus fruit</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/nectarines-plums-apricots-peaches">Nectarines, plums, apricots & peaches</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/melon-pineapple-kiwi">Kiwi & pineapple</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/melon---mango">Melon & mango</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/mango-exotic-fruit-dates-nuts">Papaya, Pomegranate & Exotic Fruit</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/dates--nuts---figs">Dates, Nuts & Figs</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/ready-to-eat">Ready to eat fruit</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/organic-fruit">Organic fruit</a></li>

                            <li><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/fresh-fruit-vegetables-special-offers">Special offers</a></li>

                    </ul>
                </div>
            </li>

    <li class="third"><a href="http://www.sainsburys.co.uk/shop/gb/groceries/fruit-veg/ripe---ready"><span>Ripe & ready</span></a>

    </li>

                    </ul>
                  </div>
                <!-- BEGIN MessageDisplay.jspf --><!-- END MessageDisplay.jspf --><!-- BEGIN ShelfDisplay.jsp -->

<div class="section">

    <h1 id="resultsHeading" class="resultsHeading">
        Ripe & ready&nbsp;(7 products available)
    </h1>

    <!-- DEBUG: shelfTopLeftESpotName = Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Left --><!-- DEBUG: shelfTopRightESpotName = Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Right -->
    <div class="eSpotContainer">

    <div id="sitecatalyst_ESPOT_NAME_Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Left" class="siteCatalystTag">Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Left</div>

    <div id="sitecatalyst_ESPOT_NAME_Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Right" class="siteCatalystTag">Z:FRUIT_AND_VEG/D:FRESH_FRUIT/A:RIPE_AND_READY/Shelf_Top_Right</div>

</div>
</div>

        <div class="section" id="filterContainer">
            <!-- FILTER SECTION STARTS HERE--><!-- BEGIN BrowseFacetsDisplay.jspf--><!-- Start Filter -->
	    <h2 class="access">Product filter options</h2>
        <div class="filterSlither">
            <div class="filterCollapseBar">
                <div class="noFlexComponent">
	                <a href="#filterOptions" id="showHideFilterSlither" aria-controls="filterOptions">
		                Filter your list
	                </a>
	                <span class="quantitySelected" id="quantitySelected" role="status" aria-live="assertive" aria-relevant="text">

	                </span>
	                <a href="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?pageSize=20&catalogId=10122&orderBy=FAVOURITES_FIRST&facet=&top_category=12518&parent_category_rn=12518&beginIndex=0&categoryId=185749&langId=44&storeId=10151" class="repressive">
	                    Clear filters
	                </a>
	            </div>
            </div>



			<form class="shelfFilterOptions " id="filterOptions" name="search_facets_form" action="" method="get" class="noFlexComponent">
	            <input type="hidden" value="44" name="langId">
	            <input type="hidden" value="10151" name="storeId">
	            <input type="hidden" value="10122" name="catalogId">
	            <input type="hidden" value="185749" name="categoryId">
	            <input type="hidden" value="12518" name="parent_category_rn">
	            <input type="hidden" value="12518" name="top_category">
	            <input type="hidden" value="20" name="pageSize">
                <input type="hidden" value="FAVOURITES_FIRST" name="orderBy">
                <input type="hidden" value="" name="searchTerm">
	            <input type="hidden" value="0" name="beginIndex">


                <div class="wrapper" id="filterOptionsContainer">



<div class="field options">
    <div class="indicator">
        <p>Options:</p>
    </div>
    <div class="checkboxes">



            <div class="input">


		                  <input id="globalOptions0" name="facet" type="checkbox" disabled="disabled" value="" aria-disabled="true" />

	                    <label class="favouritesLabel" for="globalOptions0">Favourites</label>

	        </div>



            <div class="input">


		                  <input id="globalOptions1" name="facet" type="checkbox" value="86" />

    	                <label for="globalOptions1">New</label>

	        </div>



            <div class="input">


		                  <input id="globalOptions2" name="facet" type="checkbox" disabled="disabled" value="" aria-disabled="true" />

                        <label class="offersLabel" for="globalOptions2">Offers</label>

	        </div>



    </div>

</div><!-- BEGIN BrandFacetDisplay.jspf -->

<div class="field topBrands">
    <div class="indicator">
        <p>Top Brands:</p>
    </div>
    <div class="checkboxes">


            <div class="input">

                       <input id="topBrands0" name="facet" type="checkbox" value="887" />

	           <label for="topBrands0">Sainsbury&#039;s</label>
	       </div>


    </div>
</div>

<!-- END BrandFacetDisplay.jspf -->
                </div>

                <!-- BEGIN DietaryFacetDisplay.jspf -->

<div class="filterCollapseBarDietAndLifestyle">
    <a href="#dietAndLifestyle" id="showHideDietAndLifestyle">Dietary & lifestyle options</a>
    <span class="misc">
        (such as vegetarian, organic and British)
    </span>
</div>

<div class="field dietAndLifestyle jsHide" id="dietAndLifestyle">
    <div class="checkboxes">


            <div class="input">

                        <input id="dietAndLifestyle0" name="facet" type="checkbox" value="4294966755" />

                <label for="dietAndLifestyle0">
                   Keep Refrigerated
                </label>
            </div>

            <div class="input">

                        <input id="dietAndLifestyle1" name="facet" type="checkbox" value="4294966711" />

                <label for="dietAndLifestyle1">
                   Organic
                </label>
            </div>


    </div>
</div>

<!-- END DietaryFacetDisplay.jspf -->

                <div class="filterActions">
                    <input class="button primary" type="submit" id="applyfilter" name="applyfilter" value="Apply filter" />
                </div>
            </form>
        </div>

	<!-- END BrowseFacetsDisplay.jspf--><!-- FILTER SECTION ENDS HERE-->
        </div>
        <div class="section" id="productsContainer">
            <div id="productsOverlay" class="areaOverlay"></div>
            <div id="productLister">

        <h2 class="access">Product pagination</h2>
        <div class="pagination">


    <ul class="viewOptions">

                <li class="grid">
                    <a href="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?listView=false&amp;orderBy=FAVOURITES_FIRST&amp;parent_category_rn=12518&amp;top_category=12518&amp;langId=44&amp;beginIndex=0&amp;pageSize=30&amp;catalogId=10122&amp;searchTerm=&amp;categoryId=185749&amp;listId=&amp;storeId=10151&amp;promotionId=">
                        <span class="access">Grid view</span>
                    </a>
                </li>
                <li class="listSelected">
                    <span class="access">List view</span>
                </li>

    </ul>


    <form name="search_orderBy_form" action="CategoryDisplay" method="get">
        <input type="hidden" value="44" name="langId">
        <input type="hidden" value="10151" name="storeId">
        <input type="hidden" value="10122" name="catalogId">
        <input type="hidden" value="185749" name="categoryId">
        <input type="hidden" value="20" name="pageSize">
        <input type="hidden" value="0" name="beginIndex">

        <input type="hidden" value="" name="promotionId">

        <input type="hidden" value="" name="listId">
        <input type="hidden" value="" name="searchTerm">
        <input type="hidden" name="hasPreviousOrder" value="">
        <input type="hidden" name="previousOrderId" value="" />
        <input type="hidden" name="categoryFacetId1" value="" />
        <input type="hidden" name="categoryFacetId2" value="" />
        <input type="hidden" name="bundleId" value="" />



        <div class="field">
            <div class="indicator">
                <label for="orderBy">Sort by:</label>
            </div>


                    <input type="hidden" value="12518" name="parent_category_rn">
                    <input type="hidden" value="12518" name="top_category">

            <div class="input">
                <div class="selectWrapper">
	                <select id="orderBy" name="orderBy">

	                            <option value="FAVOURITES_FIRST" selected="selected">Favourites First </option>
	                            <option value="PRICE_ASC" >Price - Low to High</option>
	                            <option value="PRICE_DESC" >Price - High to Low</option>
	                            <option value="NAME_ASC" >Product Name - A - Z</option>
	                            <option value="NAME_DESC" >Product Name - Z - A</option>
	                            <option value="TOP_SELLERS" >Top Sellers</option>

	                                <option value="RATINGS_DESC" >Ratings - High to Low</option>

	                </select>
	                <span></span>
	             </div>
            </div>
        </div>
        <div class="actions">
            <input type="submit" class="button" id="Sort" name="Sort" value="Sort" />
        </div>
    </form>


    <form name="search_pageSize_form" action="CategoryDisplay" method="get">
        <input type="hidden" value="44" name="langId">
        <input type="hidden" value="10151" name="storeId">
        <input type="hidden" value="10122" name="catalogId">
        <input type="hidden" value="185749" name="categoryId">
        <input type="hidden" value="FAVOURITES_FIRST" name="orderBy">
        <input type="hidden" value="0" name="beginIndex">

        <input type="hidden" value="" name="promotionId">
        <input type="hidden" value="" name="listId">
        <input type="hidden" value="" name="searchTerm">
        <input type="hidden" name="hasPreviousOrder" value="">
        <input type="hidden" name="previousOrderId" value="" />
        <input type="hidden" name="categoryFacetId1" value="" />
        <input type="hidden" name="categoryFacetId2" value="" />
        <input type="hidden" name="bundleId" value="" />

                <input type="hidden" value="12518" name="parent_category_rn">
                <input type="hidden" value="12518" name="top_category">

        <div class="field">
          <div class="indicator">
            <label for="pageSize">Per page</label>
          </div>
          <div class="input">
            <div class="selectWrapper">
	            <select id="pageSize" name="pageSize">

	                            <option value="20" selected="selected">20</option>

	                            <option value="40" >40</option>

	                            <option value="60" >60</option>

	                            <option value="80" >80</option>

	                            <option value="100" >100</option>

	            </select>
	            <span></span>
	         </div>
          </div>
          </div>
          <div class="actions">
              <input type="submit" class="button" id="Go" name="Go" value="Go" />
          </div>
    </form>


    <ul class="pages">
            <li class="previous">

		        <span class="access">Go to previous page</span>

            </li>

        <li class="current"><span class="access">Current page </span><span>1</span></li>

            <li class="next">

        <span class="access">Go to next page</span>

            </li>
    </ul>

       </div>

                <h2 class="access">Products</h2>
	            <ul class="productLister listView">


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error149117"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-apricot-ripe---ready-320g.html" >
	                                        Sainsbury's Apricot Ripe & Ready x5
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/wcassets/product_images/media_7572754_M.jpg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_149117">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound3.50<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound0.70<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_149117" action="OrderItemAdd" method="post" id="OrderItemAddForm_149117" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="7572754"/>

        <label class="access" for="quantity_149116">Quantity</label>

	        <input name="quantity" id="quantity_149116" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="149117"/>
        <input type="hidden" name="productId" value="149116"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_149117" id="numberInTrolley_149117">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_149117" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error572163"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado-xl-pinkerton-loose-300g.html" >
	                                        Sainsbury's Avocado Ripe & Ready XL Loose 300g
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/ExtendedSitesCatalogAssetStore/images/catalog/productImages/51/0000000202251/0000000202251_M.jpeg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_572163">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound1.50<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound1.50<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_572163" action="OrderItemAdd" method="post" id="OrderItemAddForm_572163" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="7678882"/>

        <label class="access" for="quantity_572162">Quantity</label>

	        <input name="quantity" id="quantity_572162" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="572163"/>
        <input type="hidden" name="productId" value="572162"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_572163" id="numberInTrolley_572163">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_572163" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp -->
    <div class="coverage ranged"></div>
<!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error138041"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado--ripe---ready-x2.html" >
	                                        Sainsbury's Avocado, Ripe & Ready x2
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/ExtendedSitesCatalogAssetStore/images/catalog/productImages/22/0000001600322/0000001600322_M.jpeg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_138041">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound1.80<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound1.80<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_138041" action="OrderItemAdd" method="post" id="OrderItemAddForm_138041" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="6834746"/>

        <label class="access" for="quantity_138040">Quantity</label>

	        <input name="quantity" id="quantity_138040" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="138041"/>
        <input type="hidden" name="productId" value="138040"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_138041" id="numberInTrolley_138041">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_138041" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error809817"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocados--ripe---ready-x4.html" >
	                                        Sainsbury's Avocados, Ripe & Ready x4
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/ExtendedSitesCatalogAssetStore/images/catalog/productImages/15/0000000184915/0000000184915_M.jpeg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_809817">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound3.20<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound3.20<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_809817" action="OrderItemAdd" method="post" id="OrderItemAddForm_809817" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="7718228"/>

        <label class="access" for="quantity_809816">Quantity</label>

	        <input name="quantity" id="quantity_809816" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="809817"/>
        <input type="hidden" name="productId" value="809816"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_809817" id="numberInTrolley_809817">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_809817" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error136679"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-conference-pears--ripe---ready-x4-%28minimum%29.html" >
	                                        Sainsbury's Conference Pears, Ripe & Ready x4 (minimum)
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/ExtendedSitesCatalogAssetStore/images/catalog/productImages/08/0000001514308/0000001514308_M.jpeg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_136679">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound1.50<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound1.50<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_136679" action="OrderItemAdd" method="post" id="OrderItemAddForm_136679" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="6621757"/>

        <label class="access" for="quantity_136678">Quantity</label>

	        <input name="quantity" id="quantity_136678" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="136679"/>
        <input type="hidden" name="productId" value="136678"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_136679" id="numberInTrolley_136679">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_136679" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error642875"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-golden-kiwi--taste-the-difference-x4-685641-p-44.html" >
	                                        Sainsbury's Golden Kiwi x4
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/ExtendedSitesCatalogAssetStore/images/catalog/productImages/41/0000000685641/0000000685641_M.jpeg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_642875">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound1.80<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound0.45<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_642875" action="OrderItemAdd" method="post" id="OrderItemAddForm_642875" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="685641"/>

        <label class="access" for="quantity_642874">Quantity</label>

	        <input name="quantity" id="quantity_642874" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="642875"/>
        <input type="hidden" name="productId" value="642874"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_642875" id="numberInTrolley_642875">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_642875" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>


	                            <li>
	                                <!-- BEGIN CatalogEntryThumbnailDisplay.jsp --><!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
	        <div class="errorBanner hidden" id="error130231"></div>

	        <div class="product ">
	            <div class="productInner">
	                <div class="productInfoWrapper">
	                    <div class="productInfo">

	                                <h3>
	                                    <a href="http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-kiwi-fruit--ripe---ready-x4.html" >
	                                        Sainsbury's Kiwi Fruit, Ripe & Ready x4
	                                        <img src="http://c2.sainsburys.co.uk/wcsstore7.11.1.161/SainsburysStorefrontAssetStore/wcassets/product_images/media_1116748_M.jpg" alt="" />
	                                    </a>
	                                </h3>

								<div class="ThumbnailRoundel">
								<!--ThumbnailRoundel -->
								</div>
								<div class="promoBages">
									<!-- PROMOTION -->
								</div>


	                        <!-- Review --><!-- BEGIN CatalogEntryRatingsReviewsInfo.jspf --><!-- productAllowedRatingsAndReviews: false --><!-- END CatalogEntryRatingsReviewsInfo.jspf -->
	                    </div>
	                </div>

	                <div class="addToTrolleytabBox">
	                <!-- addToTrolleytabBox LIST VIEW--><!-- Start UserSubscribedOrNot.jspf --><!-- Start UserSubscribedOrNot.jsp --><!--
			If the user is not logged in, render this opening
			DIV adding an addtional class to fix the border top which is removed
			and replaced by the tabs
		-->
		<div class="addToTrolleytabContainer addItemBorderTop">
	<!-- End AddToSubscriptionList.jsp --><!-- End AddSubscriptionList.jspf --><!--
	                        ATTENTION!!!
	                        <div class="addToTrolleytabContainer">
	                        This opening div is inside "../../ReusableObjects/UserSubscribedOrNot.jsp"
	                        -->
	                	<div class="pricingAndTrolleyOptions">
	    	                <div class="priceTab activeContainer priceTabContainer" id="addItem_130231">
	    	                    <div class="pricing">


<p class="pricePerUnit">
&pound1.80<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr>
</p>

    <p class="pricePerMeasure">&pound0.45<abbr
            title="per">/</abbr><abbr
            title="each"><span class="pricePerMeasureMeasure">ea</span></abbr>
    </p>


	    	                    </div>

	    	                                <div class="addToTrolleyForm ">

<form class="addToTrolleyForm" name="OrderItemAddForm_130231" action="OrderItemAdd" method="post" id="OrderItemAddForm_130231" class="addToTrolleyForm">
    <input type="hidden" name="storeId" value="10151"/>
    <input type="hidden" name="langId" value="44"/>
    <input type="hidden" name="catalogId" value="10122"/>
    <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
    <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
    <input type="hidden" name="SKU_ID" value="1116748"/>

        <label class="access" for="quantity_130230">Quantity</label>

	        <input name="quantity" id="quantity_130230" type="text" size="3" value="1" class="quantity"   />


        <input type="hidden" name="catEntryId" value="130231"/>
        <input type="hidden" name="productId" value="130230"/>

    <input type="hidden" name="page" value=""/>
    <input type="hidden" name="contractId" value=""/>
    <input type="hidden" name="calculateOrder" value="1"/>
    <input type="hidden" name="calculationUsage" value="-1,-2,-3"/>
    <input type="hidden" name="updateable" value="1"/>
    <input type="hidden" name="merge" value="***"/>

   	<input type="hidden" name="callAjax" value="false"/>

         <input class="button process" type="submit" name="Add" value="Add" />

</form>

	    <div class="numberInTrolley hidden numberInTrolley_130231" id="numberInTrolley_130231">

	    </div>

	    	                                </div>

	                        </div><!-- END priceTabContainer Add container --><!-- Subscribe container --><!-- Start AddToSubscriptionList.jspf --><!-- Start AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jsp --><!-- End AddToSubscriptionList.jspf -->
	                    </div>
	                </div>
	            </div>
            </div>
        	</div>
	        <div id="additionalItems_130231" class="additionalItems">
		    	<!-- BEGIN MerchandisingAssociationsDisplay.jsp --><!-- Start - JSP File Name:  MerchandisingAssociationsDisplay.jsp --><!-- END MerchandisingAssociationsDisplay.jsp -->
		    </div>

	    <!-- END CatalogEntryThumbnailDisplay.jsp -->
	                            </li>



	            </ul>


<h2 class="access">Product pagination</h2>
<div class="pagination paginationBottom">


    <ul class="pages">
            <li class="previous">

		        <span class="access">Go to previous page</span>

            </li>

        <li class="current"><span class="access">Current page </span><span>1</span></li>

            <li class="next">

        <span class="access">Go to next page</span>

            </li>
    </ul>

</div>

            </div>
        </div>
    <!-- END ShelfDisplay.jsp --><!-- ********************* ZDAS ESpot Display Start ********************** -->
            <div class="section eSpotContainer bottomESpots">
                <!-- Left POD ESpot Name = Z:FRUIT_AND_VEG/Espot_Left -->
                <!-- START ZDASPODDisplay.jsp -->

<div id="sitecatalyst_ESPOT_NAME_Z:FRUIT_AND_VEG/Espot_Left" class="siteCatalystTag">Z:FRUIT_AND_VEG/Espot_Left</div>

<div class="es es-border-box" style="width: 100%; height: 150px; padding-bottom: 15px; "><div suspendOnHover="true" dojotype="dojox.widget.AutoRotator" transition="dojox.widget.rotator.crossFade" id="myAutoRotator1445437214786" class="es-border-box-100  es-transparent-bg" duration=""><div class="es-border-box-100"><div class="es-border-box-100"><a href="/shop/gb/groceries/find-recipes/recipes/chicken-poultry-and-game/chicken--pea-and-leek-pie"><img src="http://www.sainsburys.co.uk/wcassets/2015_2016/cycle_13_18_nov/produce_recipe_leek_pie_847x135.jpg" alt="Recipe with all ingredients available to add to basket"/></a></div><div class="es-border-box es-white-bg" style="width: 168px; height: 155px; position: absolute; left: 0px; top: 0px; opacity: 1; -ms-filter: progid:DXImageTransform.Microsoft.Alpha(100); filter: alpha(opacity=100);"></div><div class="es-border-box" style="width: 168px; height: 155px; position: absolute; left: 0px; top: 0px; padding-left: 15px; padding-top: 10px; padding-right: 15px; padding-bottom: 15px; "><div class="es-border-box" style="width: 100%; padding-top: px; "><h3>Chicken, pea and leek pie</h3></div><div style="width: 100%; padding-top: px; " class="es-border-box"><p>Fluffy potato tops a creamy chicken and veg filling</p></div></div></div></div></div>
<!-- end of if empty marketingSpotDatas loop--><!-- END ZDASPODDisplay.jsp --><!--  Middle POD Espot Name = Z_Default_Espot_Content -->
                      <!-- START ZDASPODDisplay.jsp -->

<div id="sitecatalyst_ESPOT_NAME_Z_Default_Espot_Content" class="siteCatalystTag">Z_Default_Espot_Content</div>

<!-- end of if empty marketingSpotDatas loop--><!-- END ZDASPODDisplay.jsp --><!--  Right POD Espot Name = Z_Default_Espot_Content-->
                      <!-- START ZDASPODDisplay.jsp -->

<div id="sitecatalyst_ESPOT_NAME_Z_Default_Espot_Content" class="siteCatalystTag">Z_Default_Espot_Content</div>

<!-- end of if empty marketingSpotDatas loop--><!-- END ZDASPODDisplay.jsp -->
            </div>
            <!-- ********************* ZDAS ESpot Display End ********************** -->
        </div>
        <!-- content End --><!-- auxiliary Start -->
        <div class="aside" id="auxiliary">
            <!-- BEGIN RightHandSide.jspf -->
<div id="auxiliaryDock">
    <!-- BEGIN RightHandSide.jsp -->

<div class="panel loginPanel">

    <div id="sitecatalyst_ESPOT_NAME_NZ_Welcome_Back_RHS_Espot" class="siteCatalystTag">NZ_Welcome_Back_RHS_Espot</div>


	<h2>Already a customer?</h2>
    <form name="signIn" method="post" action="LogonView" id="Rhs_signIn">
        <input type="hidden" name="storeId" value="10151"/>
        <input type="hidden" name="langId" value="44"/>
        <input type="hidden" name="catalogId" value="10122"/>
        <input type="hidden" name="URL" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&categoryId=185749&langId=44&storeId=10151&krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>
        <input type="hidden" name="logonCallerId" value="LogonButton"/>
        <input type="hidden" name="errorViewName" value="CategoryDisplayView"/>
        <input class="button process" type="submit" value="Log in" />
    </form>

	<div class="panelFooter">
		<p class="register">Not yet registered?
		<a class="callToAction" name="register" href="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/PostcodeCheckView?catalogId=10122&currentPageUrl=http%3A%2F%2Fwww.sainsburys.co.uk%2Fwebapp%2Fwcs%2Fstores%2Fservlet%2FCategoryDisplay%3Fmsg%3D%26categoryId%3D185749%26langId%3D44%26storeId%3D10151%26krypto%3DdwlvaeB6%252FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%250A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%252F%252BHeNnUqybiZXu%252FL47P9A658zhrWd08mA5Y%250Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%252BardwWtMA49XQA4Iqwf%252BSvFr8RJOHK%250Afp2%252Fk0F6LH6%252Fmq5%252FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%252FydqEDvTdI5qgO6sKl0Q%253D&langId=44&storeId=10151"> Register Now</a></p>
	</div>
</div>
<div class="panel imagePanel checkPostCodePanel" id="checkPostCodePanel">

    <div id="sitecatalyst_ESPOT_NAME_NZ_Do_We_Deliver_To_You_Espot" class="siteCatalystTag">NZ_Do_We_Deliver_To_You_Espot</div>

	<h2>New customer?</h2>
    <p>Enter your postcode to check we deliver in your area.</p>


      <div id="PostCodeMessageArea" class="errorMessage" style="display:none;">
      </div>

	<form name="CheckPostCode" method="post" action="/webapp/wcs/stores/servlet/CheckPostCode" id="Rhs_checkPostCode">
		<input type="hidden" name="langId" value="44"/>
		<input type="hidden" name="storeId" value="10151"/>
		<input type="hidden" name="currentPageUrl" value="http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?msg=&amp;categoryId=185749&amp;langId=44&amp;storeId=10151&amp;krypto=dwlvaeB6%2FxULwIdnZBpXIWTi8eDrMLVBDvxz1SYU6pQ4HZ0p1fQ4WzDDbX58qo25joVKwiFFlmQW%0A0wrexmT0zSs9NxHPxri6CctBDvXHKi15cZntIRJRW%2F%2BHeNnUqybiZXu%2FL47P9A658zhrWd08mA5Y%0Azhm9vwQK7oLCWKF5VeQF9UiLmiVnffGVqRM76kUBxmRLDA%2BardwWtMA49XQA4Iqwf%2BSvFr8RJOHK%0Afp2%2Fk0F6LH6%2Fmq5%2FM97LMdaXyk0YneYUccDUWQUNnbztUkimdSo%2FydqEDvTdI5qgO6sKl0Q%3D"/>

            <input type="hidden" name="currentViewName" value="CategoryDisplayView"/>

		<input type="hidden" name="messageAreaId" value="PostCodeMessageArea"/>

		<div class="field">
			<div class="indicator">
				<label class="access" for="postCode">Postcode</label>
			</div>
			<div class="input">
				<input type="text" name="postCode" id="postCode" maxlength="8" value="" />
			</div>
		</div>
		<div class="actions">
			<input class="button primary process" type="submit" value="Check postcode"/>
		</div>
	</form>
</div>
<!-- END RightHandSide.jsp --><!-- BEGIN MiniShopCartDisplay.jsp --><!-- If we get here from a generic error this service will fail so we need to catch the exception -->
		<div class="panel infoPanel">
			<span class="icon infoIcon"></span>
		   	<h2>Important Information</h2>
			<p>Alcohol promotions available to online customers serviced from our Scottish stores may differ from those shown when browsing our site. Please log in to see the full range of promotions available to you.</p>
		</div>
	<!-- END MiniShopCartDisplay.jsp -->
</div>
<!-- END RightHandSide.jspf -->
        </div>
        <!-- auxiliary End -->
    </div>
    <!-- Main Area End --><!-- Footer Start --><!-- BEGIN LayoutContainerBottom.jspf --><!-- BEGIN FooterDisplay.jspf -->


<div id="globalFooter" class="footer">
    <ul>
	<li><a href="http://www.sainsburys.co.uk/privacy">Privacy policy</a></li>
	<li><a href="http://www.sainsburys.co.uk/cookies">Cookie policy</a></li>
	<li><a href="http://www.sainsburys.co.uk/terms">Terms &amp; conditions</a></li>
	<li><a href="http://www.sainsburys.co.uk/accessibility">Accessibility</a></li>
	<li><a href="http://help.sainsburys.co.uk/" rel="external" target="_blank" title="Opens in new window">Help Centre</a></li>
	<li><a href="http://www.sainsburys.co.uk/getintouch">Contact us</a></li>
	<li><a href="/webapp/wcs/stores/servlet/DeviceOverride?deviceId=-21&langId=44&storeId=10151">Mobile</a></li>
</ul>

</div>

<!-- END FooterDisplay.jspf --><!-- END LayoutContainerBottom.jspf --><!-- Footer Start End -->
    </div>
    <!--// End #page  --><!-- Bright Tagger start -->

	<div id="sitecatalyst_ws" class="siteCatalystTag"></div>

    <script type="text/javascript">
        var brightTagStAccount = 'sp0XdVN';
    </script>
    <noscript>
        <iframe src="//s.thebrighttag.com/iframe?c=sp0XdVN" width="1" height="1" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe>
    </noscript>

<!-- Bright Tagger End -->
</body>
</html>

<!-- END CategoriesDisplay.jsp -->
`
var (
	resultPrice float64
	resultSize  string
	resultProc  []byte

// 	itemStruct   = Items{}
// 	totUnitPrice = float64(0)
// 	doc, doc2    *goquery.Document
// 	err          error
// 	usageBool    bool
// 	jsonItems    []byte
// 	size         string
// 	url          = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)
var roundTable = []struct {
	in  float64
	out int
}{
	{float64(1.3645), 1},
	{float64(1.998), 2},
}

var toFixedTable = []struct {
	inAmt  float64
	inPrec int
	out    float64
}{
	{float64(1.3654), 2, 1.37},
	{float64(1.3645), 2, 1.36},
	{float64(1.544444444444444445), 0, 2},
	{float64(1.444444444444444445), 0, 1},
}

var getPriceTable = []struct {
	in  string
	out float64
}{
	{string("£1.80/unit"), 1.80},
	{string("£113.55/unit"), 113.55},
}

var jsonDataFail = []byte(`{"results":[{"title":"pears","size":"987kb","unit_price":3.5,"description":"pears desc"},{"title":"avacado","size":"39kb","unit_price":1.99,"description":"Avocados desc"}],"total":15.1}`)
var jsonRet = []byte(`{"results":[{"title":"Sainsbury's Apricot Ripe \u0026 Ready x5","size":"38b","unit_price":3.5,"description":"Apricots"},{"title":"Sainsbury's Avocado Ripe \u0026 Ready XL Loose 300g","size":"39b","unit_price":1.5,"description":"Avocados"},{"title":"Sainsbury's Avocado, Ripe \u0026 Ready x2","size":"43b","unit_price":1.8,"description":"Avocados"},{"title":"Sainsbury's Avocados, Ripe \u0026 Ready x4","size":"39b","unit_price":3.2,"description":"Avocados"},{"title":"Sainsbury's Conference Pears, Ripe \u0026 Ready x4 (minimum)","size":"39b","unit_price":1.5,"description":"Conference"},{"title":"Sainsbury's Golden Kiwi x4","size":"39b","unit_price":1.8,"description":"Gold Kiwi"},{"title":"Sainsbury's Kiwi Fruit, Ripe \u0026 Ready x4","size":"39b","unit_price":1.8,"description":"Kiwi"}],"total":15.1}`)
var jsonEncodeTable = []struct {
	in   Items
	json []byte
	err  error
}{
	{Items{[]Item{{"pears", "38kb", 3.5, "pears desc"}, {"avacado", "39kb", 1.99, "Avocados desc"}}, 15.1}, []byte(`{"results":[{"title":"pears","size":"38kb","unit_price":3.5,"description":"pears desc"},{"title":"avacado","size":"39kb","unit_price":1.99,"description":"Avocados desc"}],"total":15.1}`), nil},
}

func TestMain(m *testing.M) {
	usage()
	result := m.Run()
	fmt.Println("Tests complete....")
	os.Exit(result)
}

func TestGetSize(t *testing.T) {

	// Run tests in Parallel - it helps with performance and keeps the flow - we got threads let's use them.
	t.Parallel()
	// Create stupped out response using html stored in VAR mockDoc - meaning tests can be run when no connection
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, mockDoc)
	}
	// Create a mocked request
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	// Build a new test recorder which mocked response will be put into
	w := httptest.NewRecorder()
	// run mocked handler
	handler(w, req)

	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	// get the size from the body of the html.
	resultSize, err = getSize(body, "kb")
	if err != nil {
		fmt.Println(err)
	}
	// HTML is 76kb
	if resultSize != "76kb" {
		t.Logf("Expected : -->%v<--  \nbut got -->%v<--\n", string("76b"), string(resultSize))
		t.Error("Invalid size but got ", resultSize)
	}

}

func TestProcess(t *testing.T) {
	// Run tests in Parallel - it helps with performance and keeps the flow - we got threads let's use them.
	t.Parallel()
	// Get an io.Reader with HTML content
	io := strings.NewReader(mockDoc)
	// Create a qoquery document to parse from
	doc, err = goquery.NewDocumentFromReader(io)
	if err != nil {
		fmt.Println(err)
	}
	// Run the process returning a byte array with the JSON doc
	resultProc, err = process(doc)
	if err != nil {
		fmt.Println(err)
	}
	// Use Equlaity to test byte arrays.
	if !bytes.Equal(resultProc, jsonRet) {
		t.Logf("Expecteing : -->%v<--  \nbut got -->%v<--\n", string(jsonRet), string(resultProc))
		t.Error("Invalid returned Jason,expected got error ", resultProc)
	}

}

func TestJSONEncoder(t *testing.T) {
	// Run tests in Parallel - it helps with performance and keeps the flow - we got threads let's use them.
	t.Parallel()
	// Use predefined test table with input and expected output
	for _, entry := range jsonEncodeTable {
		result, _ := jsonEncoder(entry.in)
		if !bytes.Equal(result, entry.json) {
			t.Logf("Expected : -->%v<--  \nbut got -->%v<--\n", string(entry.json), string(result))
			t.Error("Invalid JSONEncoding, but got ", result)
		}
		if entry.err != nil {
			t.Error("Invalid JSONEncoding, got error ", entry.err)
		}
	}
	resFail, _ := jsonEncoder(Items{[]Item{{"pears", "99kb", 3.5, "pears desc"}, {"avacado", "39kb", 1.99, "Avocados desc"}}, 15.1})
	if bytes.Equal(resFail, jsonDataFail) {
		t.Logf("Expected inequality origin : -->%v<--  \n and returned -->%v<--\n", string(jsonDataFail), string(resFail))
		t.Error("Invalid JSONEncoding, of  ->Items{[]Item{{\"pears\", \"66kb\", 3.5, \"pears desc\"}, {\"avacado\", \"39kb\", 1.99, \"Avocados desc\"}}, 15.1}<--", string(resFail))
	}

}

func TestGetPrice(t *testing.T) {
	t.Parallel()
	for _, entry := range getPriceTable {
		resultPrice, err = getPrice(entry.in)
		if resultPrice != entry.out {
			t.Logf("Expected : -->%v<--  but got -->%v<--", entry.out, resultPrice)
			t.Error("Invalid getPrice, but got ", resultPrice)
		}
	}
}

func TestRound(t *testing.T) {
	t.Parallel()
	for _, entry := range roundTable {
		result := round(entry.in)
		if result != entry.out {
			t.Logf("Expected : -->%v<--  but got -->%v<--", entry.out, result)
			t.Error("Invalid round, but got ", result)
		}
	}
}

func TestToFixed(t *testing.T) {
	t.Parallel()
	for _, entry := range toFixedTable {
		result := toFixed(entry.inAmt, entry.inPrec)
		if result != entry.out {
			t.Logf("Expected : -->%v<--  but got -->%v<--", entry.out, result)
			t.Error("Invalid round, but got ", result)
		}
	}
}
