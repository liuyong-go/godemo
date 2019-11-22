      <?php if (!empty($favourites_explore)) {?>
            <?php foreach ($xxx as $k => $v) {?>
<div class="fav-result result row  cf">
        <div class="row col-xs-12">
          <h3 class="r-name col-xs-12 col-sm-push-5 col-sm-7 "><a href="<?php echo $this->country[COUNTRY_CODE]['base_url'];?>/ss/<?php echo $v['url_name']; ?>" title="Find out More"><?php echo $v['xx']?></a></h3>
          <input class="favorite-id" type="hidden" value="<?php echo $v['id'];?>">
          <div title="Add to favorites" ruid="<?php echo $v['xxx'];?>" class="new_fav_switch" style="cursor:pointer;"></div>

        </div>
        <div class="thumb col-xs-12 col-sm-5">
          <a href="<?php echo $this->country[COUNTRY_CODE]['base_url'];?>/xxx/<?php echo $v['url_name']; ?>"><img src="<?php echo $v['xxx'];?>" alt=" Online Restaurant Reservations" class="img-responsive"></a>
          <i class="tagline"><?php echo $v['xxx']?> </i>
        </div>
        
        <div  class="r-details row col-sm-7">
          <div class="r-info col-xs-12 col-sm-7 cf">
            <ul>
                <?php if (COUNTRY_CODE == $v['country_code']) {?>
              <li>
                <span class="type col-xs-3 col-sm-12"><?php echo _('Cuisine:');?></span>
                <p class="col-xs-9 col-sm-12">
                    <?php echo $v['cuisine']; ?> 
                  <!--<span class="cus_item">Fusion</span>,
                  <span class="cus_item">Modern European</span>-->
                 </p>
              </li>
              <li>
                <span class="type col-xs-3 col-sm-12"><?php echo _('xx:');?></span>
                <p class="col-xs-9 col-sm-12"><?php echo $v['xx']; ?></p>
              </li>
                    <?php if (!empty($v['price'])) { ?>
              <li>
                <span class="type col-xs-3 col-sm-12"><?php echo _('Price:');?></span>
                <p class="col-xs-9 col-sm-12"><?php echo str_replace("\n", "<br />", htmlspecialchars($v['xxx'], ENT_QUOTES)); ?></p>
              </li>
                    <?php }?>
                <?php } else {?>
                <li>
                <span class="type col-xs-3 col-sm-12"><?php echo _('xx:');?></span>
                <p class="col-xs-9 col-sm-12"><?php echo $v['xxx']; ?></p>
              </li>
                <?php }?>
            </ul>
          </div>
                <?php if ($v['xxx'] !=4) {?>
                    <?php if ($v['xx'] !=3) {?>
              <div class="booking col-xs-12 col-sm-5">
              <div class="col-xs-12 col-sm-12"><a href="<?php echo $book_url;?>?rid=<?php echo $v['xx'];?>&source=<?php echo $this->country[xx]['xx'];?>" class="btn"><?php echo _('xx Now');?></a></div>
            </div>
                    <?php } else {?>
            <div class="type2 booking col-xs-12 col-sm-5">
              <div class="col-xs-12 col-sm-12"><a href="<?php echo $book_url;?>?rid=<?php echo $v['xx'];?>&source=<?php echo $this->country[xx]['xx'];?>" class="btn"><?php echo _('xx Now');?></a></div>
            </div>
                    <?php }
                }?>
          </div>
                <?php if ($v['xx']) {?>
          <div class=" col-xs-12 ">
            <p><?php echo sprintf(_("xx <a href='/%s' target='_blank' class='highlight-clr'>%s </a> at this xx <a href='/%s'>(?)</a>"), $this->country[xx]['xx'], $v['xx'], $this->country[xx]['xx']);?></p>
          </div></div>
                <?php }?>
        </div>
            <?php }?>
        <?php }?>
