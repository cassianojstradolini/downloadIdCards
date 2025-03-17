package data

var benefitHTMLCardsStyle = `<style>.benefits-id-card{font-size:16px;max-width:300px}.benefits-id-card *{margin:0}.benefits-id-card .vertical-layout>*+*{margin-block-start:12px}.benefits-id-card .container{background-color:#fff;border:1px solid #d1d1d1;border-radius:12px;padding:20px;box-shadow:0 4px 8px rgba(0,0,0,.1);position:relative;overflow:hidden}.card-band::before{content:"";position:absolute;top:-30px;right:-30px;width:150px;height:10px;background-color:purple;transform:rotate(45deg)}.benefits-id-card .row{display:flex}.benefits-id-card .column{display:flex;flex-direction:column}.benefits-id-card .text-left{text-align:left}.benefits-id-card .text-right{text-align:right}.benefits-id-card .text-center{text-align:center}.benefits-id-card .title h2{color:#333}.benefits-id-card .logos .logo{width:50px;height:50px;background-image:url(https://placehold.co/50x50);background-size:contain;background-repeat:no-repeat}.benefits-id-card .label{font-weight:400;color:#666}.benefits-id-card .value{font-weight:700;color:#000}.benefits-id-card .grid-row{display:grid;grid-template-columns:1fr 1fr 1fr;grid-template-rows:auto;gap:12px}.benefits-id-card .footer{display:flex;justify-content:space-between;gap:12px}.benefits-id-card p{font-size:1em}.benefits-id-card h2{font-size:2em}.benefits-id-card h3{font-size:1.3em}.benefits-id-card h4{font-size:.1em}.benefits-id-card h5{font-size:.75em}.benefits-id-card h6{font-size:.5em}</style>`

var htmlIdCardFront = `
<div>
	` + benefitHTMLCardsStyle + `
    <section class="benefits-id-card">
        <div class="vertical-layout container card-band">
            <div class="title row text-left">
                <h2>Mock ID Card</h2>
            </div>
            <div class="vertical-layout">
                <div class="grid-row text-left">
                    <div class="column">
                        <h3 class="label light">Carrier</h3>
                        <p class="value medium">25</p>
                    </div>
                    <div class="column">
                        <h3 class="label light">Plan</h3>
                        <p class="value medium">21</p>
                    </div>
                    <div class="column">
                        <h3 class="label light">Other</h3>
                        <p class="value medium">30</p>
                    </div>
                </div>
                <div class="footer">
                    <div class="column text-left">
                        <h3 class="label light">Active:</h3>
                        <p class="value medium">Mar 1, 2024 – Feb 28, 2025</p>
                    </div>
                    <div class="logos">
                        <div aria-label="league" role="img" class="logo league"/>
                    </div>
                </div>
            </div>
        </div>
    </section>
</div>
`

var htmlIdCardBack = `
<div>
	` + benefitHTMLCardsStyle + `
    <section class="benefits-id-card">
        <div class="vertical-layout container card-band">
            <div class="title row text-left">
                <h2>Mock ID Card Back</h2>
            </div>
            <div class="text-center">
                <p>This is some back of the card information. It outlines things you can and can not do with the card.</p>
            </div>
            <div class="footer">
                <div></div>
                <div class="logos">
                    <div aria-label="league" role="img" class="logo league" />
                </div>
            </div>
        </div>
    </section>
</div>
`
var htmlIdCardBoth = `
<div>
	` + benefitHTMLCardsStyle + `
    <section class="benefits-id-card">
        <div class="vertical-layout container card-band">
            <div class="title row text-left">
                <h2>Mock ID Card</h2>
            </div>
            <div class="vertical-layout">
                <div class="grid-row text-left">
                    <div class="column">
                        <h3 class="label light">Carrier</h3>
                        <p class="value medium">25</p>
                    </div>
                    <div class="column">
                        <h3 class="label light">Plan</h3>
                        <p class="value medium">21</p>
                    </div>
                    <div class="column">
                        <h3 class="label light">Other</h3>
                        <p class="value medium">30</p>
                    </div>
                </div>
                <div class="footer">
                    <div class="column text-left">
                        <h3 class="label light">Active:</h3>
                        <p class="value medium">Mar 1, 2024 – Feb 28, 2025</p>
                    </div>
                    <div class="logos">
                        <div aria-label="league" role="img" class="logo league"/>
                    </div>
                </div>
            </div>
        </div>
    </section>
	<section class="benefits-id-card">
        <div class="vertical-layout container card-band">
            <div class="title row text-left">
                <h2>Mock ID Card Back</h2>
            </div>
            <div class="text-center">
                <p>This is some back of the card information. It outlines things you can and can not do with the card.</p>
            </div>
            <div class="footer">
                <div></div>
                <div class="logos">
                    <div aria-label="league" role="img" class="logo league" />
                </div>
            </div>
        </div>
    </section>
</div>
`
